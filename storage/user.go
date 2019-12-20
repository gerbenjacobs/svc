package storage

import (
	"context"
	"database/sql"
	"fmt"

	app "github.com/gerbenjacobs/svc"

	"github.com/google/uuid"
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
	stmt, err := u.db.PrepareContext(ctx, "INSERT INTO users (id, name, token, createdAt, updatedAt) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	uid, _ := user.ID.MarshalBinary()
	_, err = stmt.ExecContext(ctx, uid, user.Name, user.Token, user.CreatedAt, user.UpdatedAt)
	return err
}

func (u *UserRepository) Read(ctx context.Context, userID uuid.UUID) (*app.User, error) {
	uid, _ := userID.MarshalBinary()
	row := u.db.QueryRowContext(ctx, "SELECT id, name, token, createdAt, updatedAt FROM users WHERE id = ?", uid)

	user := new(app.User)
	err := row.Scan(&user.ID, &user.Name, &user.Token, &user.CreatedAt, &user.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, app.ErrUserNotFound
	case err != nil:
		return nil, fmt.Errorf("unknown error while scanning user: %v", err)
	}

	return user, nil
}
