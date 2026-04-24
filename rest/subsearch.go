package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
)

func postSubSearch(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	query := r.PostFormValue("query")
	if query == "" {
		err := fmt.Errorf("query not set")
		slog.Error(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	// XXX min length or other checks?

	result, err := model.SubscriberSearch(r.Context(), query)
	if err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, result)
}
