package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	app "github.com/gerbenjacobs/svc"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

const userCreationRequestLimit = 2048

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bytes, err := ioutil.ReadAll(io.LimitReader(r.Body, userCreationRequestLimit+1))
	if err != nil {
		error500(w, err)
		return
	}

	// bytes have reached the limit + 1, request body too big for limitreader
	if len(bytes) == userCreationRequestLimit+1 {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, handlerError{
			Code:    124,
			Message: fmt.Sprintf("Request body too large, please use %d bytes only.", userCreationRequestLimit),
		})
		return
	}

	// custom temporary User struct, with name only
	var requestBody = struct {
		Name string `json:"name"`
	}{}
	if err := json.Unmarshal(bytes, &requestBody); err != nil {
		error500(w, err)
		return
	}
	if len(requestBody.Name) > 50 {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, handlerError{
			Code:    412,
			Message: "Name has a maximum limit of 50 characters",
		})
		return
	}

	u := &app.User{
		Name:      requestBody.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	err = h.UserSvc.Add(r.Context(), u)
	if err != nil {
		error500(w, err)
		return
	}
	logrus.Infof("user created: %v", u.ID)
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, u)
}

func (h *Handler) readUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, ok := r.Context().Value(CtxKeyUserID).(string)
	if !ok {
		error500(w, errors.New("userID not found in token"))
		return
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		error500(w, errors.New("invalid user id"))
		return
	}

	user, err := h.UserSvc.User(r.Context(), uid)
	switch {
	case err == app.ErrUserNotFound:
		http.Error(w, app.ErrUserNotFound.Error(), http.StatusNotFound)
		return
	case err != nil:
		error500(w, err)
		return
	}

	logrus.WithField("user", user).Infof("user fetched: %v", user.ID)
	writeJSON(w, user)
}
