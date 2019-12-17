package storage

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"log"

	app "github.com/gerbenjacobs/svc"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(ctx context.Context, user *app.User) error {
	stmt, err := u.db.Prepare("INSERT INTO users (id, name, token, createdAt, updatedAt) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	uid, _ := user.ID.MarshalBinary()
	_, err = stmt.Exec(uid, user.Name, user.Token, user.CreatedAt, user.UpdatedAt)
	return err
}

func (u *UserRepository) Read(ctx context.Context, userID uuid.UUID) (*app.User, error) {
	uid, _ := userID.MarshalBinary()
	rows, err := u.db.Query("SELECT id, name, token, createdAt, updatedAt FROM users WHERE id = ?", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	user := new(app.User)
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Token, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}
