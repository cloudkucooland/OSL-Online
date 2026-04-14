package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getMemberChangelog(w http.ResponseWriter, r *http.Request) {
	targetIDStr := r.PathValue("id")
	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	mid := model.MemberID(targetID)

	cl, err := mid.Changelog()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cl); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}
