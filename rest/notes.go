package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getNotes(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	mid := model.MemberID(id)
	notes, err := mid.GetNotes(r.Context())
	if err != nil {
		slog.Error("failed to get notes", "member", id, "err", err)
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, notes)
}

func postNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	noteContent := r.PostFormValue("note")
	if noteContent == "" {
		sendError(w, fmt.Errorf("note content empty"), http.StatusNotAcceptable)
		return
	}

	changer, err := model.IDFromContext(r.Context())
	if err != nil {
		sendError(w, err, http.StatusUnauthorized)
		return
	}

	note := model.Note{
		Member: model.MemberID(id),
		Note:   noteContent,
		// Author: changer, // Add this if your model supports it
	}

	if err := note.Store(r.Context()); err != nil {
		slog.Error("failed to store note", "member", id, "changer", changer, "err", err)
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	noteID, err := parseID(r, "noteid")
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	n := model.NoteID(noteID)
	if err := n.Delete(r.Context()); err != nil {
		slog.Error("failed to delete note", "noteid", noteID, "err", err)
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
