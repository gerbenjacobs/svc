package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type handlerError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func error500(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func writeJSON(w http.ResponseWriter, res interface{}) {
	e := json.NewEncoder(w)
	e.SetIndent("", " ")
	_ = e.Encode(res)
}

func health(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprint(w, "OK!")
}

func redirect(url string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}
