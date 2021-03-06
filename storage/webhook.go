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

func (w *WebhookRepository) Create(ctx context.Context, webhook *app.Webhook) error {
	return nil
}

func (w *WebhookRepository) Read(ctx context.Context, webhookID string) (*app.Webhook, error) {
	return nil, fmt.Errorf("%q: %w", webhookID, app.ErrWebhookNotFound)
}

func (w *WebhookRepository) Update(ctx context.Context, webhook *app.Webhook) error {
	return nil
}

func (w *WebhookRepository) Delete(ctx context.Context, webhookID string) error {
	return nil
}
