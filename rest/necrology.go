package rest

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getNecrology(w http.ResponseWriter, r *http.Request) {
	slog.Info("necrology", "requester", getUser(r))
	isee, err := model.Necrology(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, isee)
}

func getCommemorations(w http.ResponseWriter, r *http.Request) {
	date := time.Now()
	month := date.Month()
	day := date.Day()

	if m := r.FormValue("month"); m != "" {
		mm, err := parseIDFromString(m)
		if err != nil {
			slog.Error(err.Error())
			sendError(w, err, http.StatusInternalServerError)
			return
		}
		month = time.Month(mm)
	}

	if d := r.FormValue("day"); d != "" {
		dd, err := parseIDFromString(d)
		if err != nil {
			slog.Error(err.Error())
			sendError(w, err, http.StatusInternalServerError)
			return
		}
		day = dd
	}

	slog.Info("getCommemorations", "month", month, "day", day)
	isee, err := model.Commemorations(r.Context(), month, day)
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, isee)
}
