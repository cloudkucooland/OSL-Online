package rest

import (
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getLeadership(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")
	if category == "" {
		category = "elected"
	}

	slog.Info("leadership", "category", category, "requester", getUser(r))
	leaders, err := model.Leadership(r.Context(), category)
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, leaders)
}
