package services

import (
	"context"

	"github.com/google/uuid"

	app "github.com/gerbenjacobs/svc"
)

// Rationale: Ideally these methods say what they do according to business value
// in this example however it's really just CRUD-only.
// However imagine bigger services having methods like FindNewUsers or GetLoyalCustomers
type UserService interface {
	Add(ctx context.Context, user *app.User) error
	User(ctx context.Context, userID uuid.UUID) (*app.User, error)
}

// Rationale: It's tempting to name your methods like they are below, but this is actually a bit of a smell
// it will make you wonder why you need this interface if it looks a lot like storage.WebhookStorage.
// I suggest you to keep the service and storage interface separated.
type WebhookService interface {
	Create(ctx context.Context, webhook *app.Webhook) error
	Read(ctx context.Context, webhookID string) (*app.Webhook, error)
	Update(ctx context.Context, webhook *app.Webhook) error
	Delete(ctx context.Context, webhookID string) error
}
