package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func getLocalities(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	l, err := model.Localities()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(l); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func getLocalityMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	q := ps.ByName("joint")

	slog.Info("locality search", "locality", q)

	chunks := strings.Split(q, "-")
	var country, locality string
	country = chunks[0]
	if len(chunks) > 1 {
		locality = chunks[1]
	}

	m, err := model.LocalityMembers(country, locality)
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
