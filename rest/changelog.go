package rest

import (
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getMemberChangelog(w http.ResponseWriter, r *http.Request) {
	targetID, err := parseID(r, "id")
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	mid := model.MemberID(targetID)

	cl, err := mid.Changelog(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, cl)
}
