package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func getNotes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	s := ps.ByName("id")
	id, err := strconv.Atoi(s)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	m := model.MemberID(id)
	slog.Info("get notes", "member", id, "requester", getUser(r))
	notes, err := m.GetNotes()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func postNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	s := ps.ByName("id")
	if s == "" {
		err := fmt.Errorf("id not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}
	id, err := strconv.Atoi(s)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}
	m := model.MemberID(id)

	notecontent := r.PostFormValue("note")
	if notecontent == "" {
		err := fmt.Errorf("note not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	note := model.Note{
		ID:     0,
		Member: m,
		Note:   notecontent,
	}
	if err := note.Store(); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func deleteNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	s := ps.ByName("id")
	if s == "" {
		err := fmt.Errorf("id not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}
	id, err := strconv.Atoi(s)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}
	n := model.NoteID(id)

	if err := n.Delete(); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
