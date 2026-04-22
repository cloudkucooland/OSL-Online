package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"
)

func postSearch(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	query := r.PostFormValue("query")
	if query == "" || query == "undefined" {
		err := fmt.Errorf("query not set")
		slog.Error(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	query = strings.TrimSpace(query)
	if len(query) < 3 {
		err := fmt.Errorf("query too short")
		slog.Error(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	slog.Info("search", "query", query, "requester", getUser(r))
	result, err := model.Search(r.Context(), query)
	if err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, result)
}

func postEmailSearch(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	query := r.PostFormValue("query")
	if query == "" || query == "undefined" {
		err := fmt.Errorf("query not set")
		slog.Error(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	slog.Info("search email", "query", query, "requester", getUser(r))
	result, err := model.SearchEmail(r.Context(), query)
	if err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, result)
}
