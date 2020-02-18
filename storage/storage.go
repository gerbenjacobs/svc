package storage

import (
	"context"

	"github.com/google/uuid"

	app "github.com/gerbenjacobs/svc"
)

type UserStorage interface {
	Create(ctx context.Context, user *app.User) error
	Read(ctx context.Context, userID uuid.UUID) (*app.User, error)
	// Rationale: These methods are not implemented in this example.
	// Update(ctx context.Context, user *app.User) error
	// Delete(ctx context.Context, userID string) error
}

type WebhookStorage interface {
	Create(ctx context.Context, webhook *app.Webhook) error
	Read(ctx context.Context, webhookID string) (*app.Webhook, error)
	Update(ctx context.Context, webhook *app.Webhook) error
	Delete(ctx context.Context, webhookID string) error
}
