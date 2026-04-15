package rest

import (
	"context"
	"encoding/json"
	"github.com/cloudkucooland/OSL-Online/model"
	"log/slog"
	"net/http"
	"strconv"
)

func getPublicPrayers(w http.ResponseWriter, r *http.Request) {
	// this is not wrapped in authMW, but we could still check for the JWT and set proper values...
	ctx := context.WithValue(r.Context(), model.CtxKeyID, model.MemberID(0))
	ctx = context.WithValue(ctx, model.CtxKeyLevel, model.AuthLevelView)

	prayers, err := model.GetPrayers(ctx, nil, true)
	if err != nil {
		slog.Error("public prayers", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(prayers)
}

func getPrayers(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusNotAcceptable)
		return
	}

	mid := model.MemberID(id)
	prayers, err := model.GetPrayers(r.Context(), &mid, false)
	if err != nil {
		slog.Error("member prayers", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(prayers)
}

func addPrayer(w http.ResponseWriter, r *http.Request) {
	var p model.Prayer
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		slog.Error("add prayer", "error", err, "incoming", p)
		http.Error(w, "bad request", http.StatusNotAcceptable)
		return
	}

	if err := p.Insert(r.Context()); err != nil {
		slog.Error("add prayer", "error", err)
		http.Error(w, "could not save", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func deletePrayer(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	if err := model.DeletePrayer(r.Context(), model.PrayerID(id)); err != nil {
		slog.Error("delete prayer", "error", err, "id", id)
		http.Error(w, "could not delete", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
