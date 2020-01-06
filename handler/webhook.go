package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	app "github.com/gerbenjacobs/svc"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) createWebhook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.Error(w, "webhook creation not implemented yet", http.StatusNotFound)
}

func (h *Handler) readWebhook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := h.WebhookSvc.Read(r.Context(), p.ByName("webhookID"))
	switch {
	case errors.Is(err, app.ErrWebhookNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case err != nil:
		error500(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		error500(w, err)
		return
	}
}

func (h *Handler) updateWebhook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.Error(w, "webhook updating not implemented yet", http.StatusNotFound)
}

func (h *Handler) deleteWebhook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if err := h.WebhookSvc.Delete(r.Context(), p.ByName("webhookID")); err != nil {
		error500(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
