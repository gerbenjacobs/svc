package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	app "github.com/gerbenjacobs/svc"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) createWebhook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.Error(w, "webhook creation not implemented yet", http.StatusNotFound)
}

func (h *Handler) readWebhook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	webhook, err := h.WebhookSvc.Read(r.Context(), p.ByName("webhookID"))
	switch {
	case errors.Is(err, app.ErrWebhookNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case err != nil:
		error500(w, err)
		return
	}

	// custom output format for webhooks
	// Rationale: this is just to simulate that it's okay to convert a domain model
	// to an expected output format for HTTP responses.
	type webhookOutput struct {
		URL         string    `json:"url"`
		Triggers    []string  `json:"triggers"`
		TriggeredAt time.Time `json:"triggered_at"`
	}

	whResp := webhookOutput{
		URL:         webhook.URL,
		Triggers:    webhook.Events,
		TriggeredAt: webhook.UpdatedAt,
	}

	if err := json.NewEncoder(w).Encode(whResp); err != nil {
		error500(w, err)
		return
	}
}

func (h *Handler) updateWebhook(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	http.Error(w, "webhook updating not implemented yet", http.StatusNotFound)
}

func (h *Handler) deleteWebhook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if err := h.WebhookSvc.Delete(r.Context(), p.ByName("webhookID")); err != nil {
		error500(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
