package svc

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrWebhookNotFound = errors.New("webhook not found")

type Webhook struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"-"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
