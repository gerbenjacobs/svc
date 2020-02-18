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

	// Rationale: I'm reusing the app.User here because the fields are quite primitive types
	// Depending on your scheme you could easily do some transformations here to change
	// app.User to a customer UserDAO struct, f.e. when your database engine stores bools as tinyints.
	user := new(app.User)
	err := row.Scan(&user.ID, &user.Name, &user.Token, &user.CreatedAt, &user.UpdatedAt)
	switch {
	// Rationale: Our service layer knows nothing about sql.ErrNoRows, but we at this point do
	// that's why it's important to convert your database engine errors to common Domain model errors
	// that are known within the application.
	// This specific example makes use of the %w verb to wrap errors with a custom message
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("user with ID %q not found: %w", userID, app.ErrUserNotFound)
	// Rationale: Here we're explicitly not wrapping the error as the service shouldn't do anything with it.
	// However if you started noticing these in your logs, you can probably handle them like the above case.
	case err != nil:
		return nil, fmt.Errorf("unknown error while scanning user: %v", err)
	}

	return user, nil
}
