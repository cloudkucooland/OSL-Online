package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getSubscriber(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusBadRequest)
		return
	}
	sid := model.SubscriberID(id)
	m, err := sid.Get(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	m.FormattedAddr, err = model.FormatAddress(m)
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	slog.Info("getSubscriber", "id", id, "requester", getUser(r))
	sendJSON(w, m)
}

func setSubscriber(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	field := r.PostFormValue("field")
	if field == "" {
		sendError(w, fmt.Errorf("field not set"), http.StatusNotAcceptable)
		return
	}

	value := r.PostFormValue("value")

	id, err := parseID(r, "id")
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusBadRequest)
		return
	}
	sid := model.SubscriberID(id)
	if err := sid.SetField(r.Context(), field, value); err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
