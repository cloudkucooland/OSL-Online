package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getNotes(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	mid := model.MemberID(id)
	notes, err := mid.GetNotes()
	if err != nil {
		slog.Error("failed to get notes", "member", id, "err", err)
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(notes)
}

func postNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	noteContent := r.PostFormValue("note")
	if noteContent == "" {
		http.Error(w, jsonError(fmt.Errorf("note content empty")), http.StatusNotAcceptable)
		return
	}

	changer, err := model.IDFromContext(r.Context())
	if err != nil {
		http.Error(w, jsonError(err), http.StatusUnauthorized)
		return
	}

	note := model.Note{
		Member: model.MemberID(id),
		Note:   noteContent,
		// Author: changer, // Add this if your model supports it
	}

	if err := note.Store(); err != nil {
		slog.Error("failed to store note", "member", id, "changer", changer, "err", err)
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	noteID, err := strconv.Atoi(r.PathValue("noteid"))
	if err != nil {
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	n := model.NoteID(noteID)
	if err := n.Delete(); err != nil {
		slog.Error("failed to delete note", "noteid", noteID, "err", err)
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
