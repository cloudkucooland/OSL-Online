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

func getMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	unlisted := false
	level, err := getLevel(r)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	if level >= AuthLevelManager {
		unlisted = true
	}

	mid := model.MemberID(id)
	m, err := mid.Get(unlisted)
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

func getMemberChapters(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	id, err := strconv.Atoi(ps.ByName("id"))
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

func setMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	email := getUser(r)
	changer, err := model.GetID(email)
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

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := model.SetMemberField(model.MemberID(id), field, value, model.MemberID(changer)); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func createMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	firstname := r.PostFormValue("firstname")
	if firstname == "" {
		err := fmt.Errorf("firstname not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	lastname := r.PostFormValue("lastname")
	if lastname == "" {
		err := fmt.Errorf("lastname not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	id, err := model.Create(firstname, lastname)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	out := fmt.Sprintf(`{"status":"ok", "id": %d}`, id)
	fmt.Fprint(w, out)
}

func setMemberChapters(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	id, err := strconv.Atoi(ps.ByName("id"))
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
