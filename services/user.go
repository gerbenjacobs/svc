package services

import (
	"context"
	"time"

	app "github.com/gerbenjacobs/svc"
	"github.com/gerbenjacobs/svc/storage"

	"github.com/google/uuid"
)

// UserSvc is our service struct that implements the services.UserService interface
type UserSvc struct {
	storage storage.UserStorage
	auth    *Auth
}

func NewUserSvc(userStorage storage.UserStorage, auth *Auth) (*UserSvc, error) {
	return &UserSvc{
		storage: userStorage,
		auth:    auth,
	}, nil
}

// User returns the user based on the user ID
func (u *UserSvc) User(ctx context.Context, userID uuid.UUID) (*app.User, error) {
	return u.storage.Read(ctx, userID)
}

// Add adds a user to our service and repository
func (u *UserSvc) Add(ctx context.Context, user *app.User) error {
	userID := uuid.New()

	token, err := u.auth.Create(userID.String())
	if err != nil {
		return err
	}

	// create user object
	user.ID = userID
	user.Token = token
	n := time.Now().UTC()
	user.CreatedAt = n
	user.UpdatedAt = n

	// persist it
	return u.storage.Create(ctx, user)
}
