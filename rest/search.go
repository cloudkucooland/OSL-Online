package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"
)

func postSearch(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	query := r.PostFormValue("query")
	if query == "" || query == "undefined" {
		err := fmt.Errorf("query not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	query = strings.TrimSpace(query)
	if len(query) < 3 {
		err := fmt.Errorf("query too short")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	slog.Info("search", "query", query, "requester", getUser(r))
	result, err := model.Search(r.Context(), query)
	if err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func postEmailSearch(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	query := r.PostFormValue("query")
	if query == "" || query == "undefined" {
		err := fmt.Errorf("query not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	slog.Info("search email", "query", query, "requester", getUser(r))
	result, err := model.SearchEmail(r.Context(), query)
	if err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}
