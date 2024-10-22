package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func getMemberChangelog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	m, err := model.GetMember(id, false)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	cl, err := m.Changelog()
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