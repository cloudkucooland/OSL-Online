package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"

	"github.com/julienschmidt/httprouter"
)

func getDashboard(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result, err := model.Dashboard()
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
