package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func getMe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	email := getUser(r)
	id, err := model.GetID(email)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	mid := model.MemberID(id)
	m, err := mid.Get(true)
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

func getMeChapters(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	email := getUser(r)
	id, err := model.GetID(email)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	mid := model.MemberID(id)
	m, err := mid.Get(true)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	chapters, err := m.GetChapters()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(chapters); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func setMe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	email := getUser(r)
	id, err := model.GetID(email)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	field := r.PostFormValue("field")
	if field == "" {
		err := fmt.Errorf("field not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	value := r.PostFormValue("value")
	if err := model.SetMeField(id, field, value); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func setMeChapters(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	email := getUser(r)
	id, err := model.GetID(email)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	cs := r.PostFormValue("chapters")
	ss := strings.Split(cs, ",")

	chapters := make([]int, 0)
	for _, n := range ss {
		if n == "" {
			continue
		}
		c, err := strconv.Atoi(n)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		chapters = append(chapters, c)
	}

	mid := model.MemberID(id)
	member, err := mid.Get(true)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := member.SetChapters(chapters...); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
