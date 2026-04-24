package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getChapters(w http.ResponseWriter, r *http.Request) {
	ch, err := model.Chapters(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	slog.Info("getChapters", "requester", getUser(r))
	sendJSON(w, ch)
}

func getChapterMembers(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusBadRequest)
		return
	}

	chapter := model.Chapter{
		ID: model.ChapterID(id),
	}
	slog.Info("chapter membership", "chapter", id, "requester", getUser(r))
	members, err := chapter.Members(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, members)
}

func putChapter(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	id, err := parseID(r, "id")
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	if name == "" {
		sendError(w, fmt.Errorf("name not set"), http.StatusNotAcceptable)
		return
	}

	priorID, err := parseIDFromString(r.PostFormValue("prior"))
	if err != nil {
		sendError(w, fmt.Errorf("prior id not set or invalid: %w", err), http.StatusNotAcceptable)
		return
	}

	c := model.Chapter{
		ID:    model.ChapterID(id),
		Name:  name,
		Prior: model.MemberID(priorID),
	}
	if err := c.Update(r.Context()); err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
