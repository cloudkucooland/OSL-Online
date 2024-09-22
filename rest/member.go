package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func getMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	m, err := model.GetMember(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}
