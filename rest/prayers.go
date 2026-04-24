package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getPublicPrayers(w http.ResponseWriter, r *http.Request) {
	// this is not wrapped in authMW, but we could still check for the JWT and set proper values...
	ctx := context.WithValue(r.Context(), model.CtxKeyID, model.MemberID(0))
	ctx = context.WithValue(ctx, model.CtxKeyLevel, model.AuthLevelView)

	prayers, err := model.GetPrayers(ctx, nil, true)
	if err != nil {
		slog.Error("public prayers", "error", err)
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	sendJSON(w, prayers)
}

func getPrayers(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	mid := model.MemberID(id)
	prayers, err := model.GetPrayers(r.Context(), &mid, false)
	if err != nil {
		slog.Error("member prayers", "error", err)
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	sendJSON(w, prayers)
}

func addPrayer(w http.ResponseWriter, r *http.Request) {
	var p model.Prayer
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		slog.Error("add prayer", "error", err, "incoming", p)
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	if err := p.Insert(r.Context()); err != nil {
		slog.Error("add prayer", "error", err)
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func deletePrayer(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	if err := model.DeletePrayer(r.Context(), model.PrayerID(id)); err != nil {
		slog.Error("delete prayer", "error", err, "id", id)
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
