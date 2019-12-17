package services

import (
	"context"

	app "github.com/gerbenjacobs/svc"
	"github.com/gerbenjacobs/svc/storage"
)

type WebhookSvc struct {
	storage storage.WebhookStorage
}

func NewWebhookService(webhookStorage storage.WebhookStorage) *WebhookSvc {
	return &WebhookSvc{
		storage: webhookStorage,
	}
}

func (w *WebhookSvc) Create(ctx context.Context, webhook *app.Webhook) error {
	return w.storage.Create(ctx, webhook)
}

func (w *WebhookSvc) Read(ctx context.Context, webhookID string) (*app.Webhook, error) {
	return w.storage.Read(ctx, webhookID)
}

func (w *WebhookSvc) Update(ctx context.Context, webhook *app.Webhook) error {
	return w.storage.Update(ctx, webhook)
}

func (w *WebhookSvc) Delete(ctx context.Context, webhookID string) error {
	return w.storage.Delete(ctx, webhookID)
}
