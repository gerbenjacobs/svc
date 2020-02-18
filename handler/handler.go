package handler

import (
	"net/http"

	"github.com/gerbenjacobs/svc/services"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Handler is your dependency container
// Rationale: This acts as your entry point for your application, everything is delegated from here
type Handler struct {
	mux http.Handler
	Dependencies
}

// Dependencies contains all the dependencies your application and its services require
// Rationale: Here you should use the interfaces of your services, although actual implementations can also be passed
// but note that this limits you while mocking for example.
type Dependencies struct {
	UserSvc    services.UserService
	WebhookSvc services.WebhookService
	Auth       *services.Auth
}

// New creates a new handler given a set of dependencies
// Rationale: This function, specifically in handler/handler.go contains all your routes.
// These are HTTP endpoints, but you can also switch it out with GRPC for example.
func New(dependencies Dependencies) *Handler {
	h := &Handler{
		Dependencies: dependencies,
	}

	r := httprouter.New()
	r.GET("/", redirect("health"))
	r.GET("/health", health)

	r.POST("/v1/user", h.createUser)
	r.GET("/v1/user", h.AuthMiddleware(h.readUser))

	// Rationale: None of these are implemented, it's just to show having several services/repositories.
	r.POST("/v1/webhook", h.createWebhook)
	r.GET("/v1/webhook/:webhookID", h.readWebhook)
	r.PUT("/v1/webhook/:webhookID", h.updateWebhook)
	r.DELETE("/v1/webhook/:webhookID", h.deleteWebhook)

	// create chained list of middleware
	// and wrap our router with it
	mw := alice.New(customLoggingMiddleware)
	h.mux = mw.Then(r)
	return h
}

// ServeHTTP makes sure Handler implements the http.Handler interface
// this keeps the underlying mux private
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
