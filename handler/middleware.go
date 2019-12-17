package handler

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

const (
	CtxKeyUserID = iota
)

func (h *Handler) AuthMiddleware(f httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		claims, err := h.Auth.ReadFromRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// Add user ID in context
		r = r.WithContext(context.WithValue(r.Context(), CtxKeyUserID, claims.UserID))
		f(w, r, p)
	}
}

func customLoggingMiddleware(handler http.Handler) http.Handler {
	return handlers.CustomLoggingHandler(os.Stdout, handler, func(_ io.Writer, p handlers.LogFormatterParams) {
		logrus.Debugf("%d %s \"%s %s\" %d \"%s\"", p.StatusCode, p.Request.Proto, p.Request.Method, p.URL.String(), p.Size, p.Request.Header.Get("User-Agent"))
	})
}
