package handler

import (
	"fmt"
	"net/http"

	"github.com/gerbenjacobs/svc/services"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type Handler struct {
	mux http.Handler
	Dependencies
}

type Dependencies struct {
	UserSvc    services.UserService
	WebhookSvc services.WebhookService
	Auth       *services.Auth
}

type handlerError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func New(dependencies Dependencies) *Handler {
	h := &Handler{
		Dependencies: dependencies,
	}

	r := httprouter.New()
	r.GET("/", redirect("health"))
	r.GET("/health", health)

	r.POST("/v1/user", h.createUser)
	r.GET("/v1/user", h.AuthMiddleware(h.readUser))

	r.POST("/v1/webhook", h.createWebhook)
	r.GET("/v1/webhook/:webhookID", h.readWebhook)
	r.PUT("/v1/webhook/:webhookID", h.updateWebhook)
	r.DELETE("/v1/webhook/:webhookID", h.deleteWebhook)

	mw := alice.New(customLoggingMiddleware)
	h.mux = mw.Then(r)
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func health(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprint(w, "OK!")
}

func redirect(url string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}
