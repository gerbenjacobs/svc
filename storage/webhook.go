package storage

import (
	"context"
	"database/sql"
	"fmt"

	app "github.com/gerbenjacobs/svc"
)

type WebhookRepository struct {
	db *sql.DB
}

func NewWebhookRepository(db *sql.DB) *WebhookRepository {
	return &WebhookRepository{
		db: db,
	}
}

func (w *WebhookRepository) Create(context.Context, *app.Webhook) error {
	return nil
}

func (w *WebhookRepository) Read(ctx context.Context, webhookID string) (*app.Webhook, error) {
	return nil, fmt.Errorf("%q: %w", webhookID, app.ErrWebhookNotFound)
}

func (w *WebhookRepository) Update(context.Context, *app.Webhook) error {
	return nil
}

func (w *WebhookRepository) Delete(context.Context, string) error {
	return nil
}
