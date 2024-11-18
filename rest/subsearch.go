package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"

	"github.com/julienschmidt/httprouter"
)

func postSubSearch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	query := r.PostFormValue("query")
	if query == "" {
		err := fmt.Errorf("query not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	// XXX min length or other checks?

	result, err := model.SubscriberSearch(query)
	if err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	headers(w, r)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}
