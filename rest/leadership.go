package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func getLeadership(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	category := ps.ByName("category")
	if category == "" {
		category = "elected"
	}

	slog.Info("leadership", "category", category, "requester", getUser(r))
	leaders, err := model.Leadership(r.Context(), category)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(leaders); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}
