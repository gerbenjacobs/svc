package services

import (
	"context"

	"github.com/google/uuid"

	app "github.com/gerbenjacobs/svc"
)

type UserService interface {
	Add(ctx context.Context, user *app.User) error
	User(ctx context.Context, userID uuid.UUID) (*app.User, error)
}

type WebhookService interface {
	Create(ctx context.Context, webhook *app.Webhook) error
	Read(ctx context.Context, webhookID string) (*app.Webhook, error)
	Update(ctx context.Context, webhook *app.Webhook) error
	Delete(ctx context.Context, webhookID string) error
}
