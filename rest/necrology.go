package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func getNecrology(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	slog.Info("necrology", "requester", getUser(r))
	isee, err := model.Necrology(r.Context())
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

func getCommemorations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	date := time.Now()
	month := date.Month()
	day := date.Day()

	if m := r.FormValue("month"); m != "" {
		mm, err := strconv.Atoi(m)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, jsonError(err), http.StatusInternalServerError)
			return
		}
		month = time.Month(mm)
	}

	if d := r.FormValue("day"); d != "" {
		dd, err := strconv.Atoi(d)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, jsonError(err), http.StatusInternalServerError)
			return
		}
		day = dd
	}

	slog.Info("getCommemorations", "month", month, "day", day)
	isee, err := model.Commemorations(r.Context(), month, day)
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
