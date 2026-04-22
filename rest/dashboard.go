package rest

import (
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getDashboard(w http.ResponseWriter, r *http.Request) {
	result, err := model.Dashboard(r.Context())
	if err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	slog.Info("getDashboard", "requester", getUser(r))
	sendJSON(w, result)
}
