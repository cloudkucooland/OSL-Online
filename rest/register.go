package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/cloudkucooland/OSL-Online/email"
	"github.com/cloudkucooland/OSL-Online/model"
)

func postRegister(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	username := model.Authname(strings.TrimSpace(r.PostFormValue("email")))
	if username == "" {
		err := fmt.Errorf("email not set")
		slog.Error(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	password, err := username.Register(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	if err := email.SendRegister(username.String(), password); err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
