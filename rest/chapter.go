package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getChapters(w http.ResponseWriter, r *http.Request) {
	ch, err := model.Chapters()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	slog.Info("getChapters", "requester", getUser(r))
	if err := json.NewEncoder(w).Encode(ch); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func getChapterMembers(w http.ResponseWriter, r *http.Request) {
	s := r.PathValue("id")
	id, err := strconv.Atoi(s)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	chapter := model.Chapter{
		ID: model.ChapterID(id),
	}
	slog.Info("chapter membership", "chapter", id, "requester", getUser(r))
	members, err := chapter.Members(r.Context())
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(members); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func putChapter(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	s := r.PathValue("id")
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

	name := r.PostFormValue("name")
	if name == "" {
		err := fmt.Errorf("name not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	s = r.PostFormValue("prior")
	if s == "" {
		err := fmt.Errorf("prior id not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}
	prior, err := strconv.Atoi(s)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	c := model.Chapter{
		ID:    model.ChapterID(id),
		Name:  name,
		Prior: model.MemberID(prior),
	}
	if err := c.Update(); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
