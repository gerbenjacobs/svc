package handler

import (
	"encoding/json"
	"net/http"
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
