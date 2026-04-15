package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getDashboard(w http.ResponseWriter, r *http.Request) {
	result, err := model.Dashboard(r.Context())
	if err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	slog.Info("getDashboard", "requester", getUser(r))
	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}
