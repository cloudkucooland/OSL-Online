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

	leaders, err := model.Leadership(category)
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
