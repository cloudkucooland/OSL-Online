package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getMe(w http.ResponseWriter, r *http.Request) {
	username := model.Authname(getUser(r))
	id, err := username.GetID(r.Context())

	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	mid := model.MemberID(id)
	m, err := mid.Get(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	slog.Info("loading self", "user", m.OSLName())
	sendJSON(w, m)
}

func getMeChapters(w http.ResponseWriter, r *http.Request) {
	username := model.Authname(getUser(r))
	id, err := username.GetID(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	mid := model.MemberID(id)
	m, err := mid.Get(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	chapters, err := m.ID.GetChapters(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, chapters)
}

func setMe(w http.ResponseWriter, r *http.Request) {
	id, err := model.IDFromContext(r.Context())
	if err != nil {
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	field := r.PostFormValue("field")
	if field == "" {
		sendError(w, fmt.Errorf("field not set"), http.StatusNotAcceptable)
		return
	}

	value := r.PostFormValue("value")
	if err := model.SetMeField(r.Context(), id, field, value); err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, jsonStatusOK)
}

func setMeChapters(w http.ResponseWriter, r *http.Request) {
	id, err := model.IDFromContext(r.Context())
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	cs := r.PostFormValue("chapters")
	ss := strings.Split(cs, ",")

	chapters := make([]int, 0)
	for _, n := range ss {
		if n == "" {
			continue
		}
		c, err := parseIDFromString(n)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		chapters = append(chapters, c)
	}

	mid := model.MemberID(id)
	member, err := mid.Get(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	if err := member.SetChapters(r.Context(), chapters...); err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func getMeGiving(w http.ResponseWriter, r *http.Request) {
	username := model.Authname(getUser(r))
	id, err := username.GetID(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	mid := model.MemberID(id)
	gr, err := mid.GivingRecords(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	slog.Info("loading self giving record", "user", username)
	sendJSON(w, gr)
}
