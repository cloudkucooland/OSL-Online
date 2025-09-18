package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func getNecrology(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	slog.Info("necrology", "requester", getUser(r))
	isee, err := model.Necrology()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(isee); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}
