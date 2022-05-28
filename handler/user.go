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
	"github.com/gerbenjacobs/svc/services"

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
			Code:    124, // this code doesn't mean anything in our API, but it could..
			Message: fmt.Sprintf("Request body too large, please use %d bytes only.", userCreationRequestLimit),
		})
		return
	}

	// custom temporary User struct, with name only
	// Rationale: In this theoretical situation I don't want users to touch my app.User object straight away,
	// so I create this auxiliary struct to temporary store the only information I want from the user.
	var requestBody = struct {
		Name string `json:"name"`
	}{}
	if err := json.Unmarshal(bytes, &requestBody); err != nil {
		error500(w, err)
		return
	}
	if len(requestBody.Name) > services.MaxUsernameLength {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, handlerError{
			Code:    412, // this code doesn't mean anything in our API, but it could..
			Message: fmt.Sprintf("Name has a maximum limit of %d characters", services.MaxUsernameLength),
		})
		return
	}

	// Rationale: Since we're about to enter the service layer, I need to start talking in Domain Models,
	// so I change my auxiliary struct into an actual app.User.
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
	logrus.Infof("user created: [%v] %s", u.ID, u.Name)
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
	case errors.Is(err, app.ErrUserNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case err != nil:
		error500(w, err)
		return
	}

	logrus.WithField("user", user).Infof("user fetched: %v", user.ID)
	writeJSON(w, user)
}
