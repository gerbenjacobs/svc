package svc

import (
	"errors"
)

// Rationale: Depending on how big your services/application will be it's up to you to decide
// whether you use an error.go file or have errors inside the files for your models.
// So far my services didn't have that many errors that made this file become unwieldy.
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrWebhookNotFound = errors.New("webhook not found")
)
