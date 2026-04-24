package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getMember(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		slog.Error("invalid id", "err", err)
		sendError(w, err, http.StatusBadRequest)
		return
	}

	mid := model.MemberID(id)
	m, err := mid.Get(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	m.FormattedAddr, err = model.FormatAddress(m)
	if err != nil {
		slog.Error("address format failed", "err", err)
	}

	slog.Info("user loaded", "user", m.OSLName(), "requester", getUser(r))
	sendJSON(w, m)
}

func getMemberChapters(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	mid := model.MemberID(id)
	chapters, err := mid.GetChapters(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, chapters)
}

func setMember(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	field := r.PostFormValue("field")
	value := r.PostFormValue("value")

	if field == "" {
		sendError(w, fmt.Errorf("field not set"), http.StatusNotAcceptable)
		return
	}

	if err := model.MemberID(id).SetMemberField(r.Context(), field, value); err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func createMember(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	firstname := r.PostFormValue("firstname")
	lastname := r.PostFormValue("lastname")

	if firstname == "" || lastname == "" {
		sendError(w, fmt.Errorf("name components missing"), http.StatusNotAcceptable)
		return
	}

	id, err := model.Create(firstname, lastname)
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `{"status":"ok", "id": %d}`, id)
}

func setMemberChapters(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}
	cs := r.PostFormValue("chapters")
	var chapters []int

	if cs != "" {
		ss := strings.Split(cs, ",")
		for _, n := range ss {
			if c, err := parseIDFromString(strings.TrimSpace(n)); err == nil {
				chapters = append(chapters, c)
			}
		}
	}

	mid := model.MemberID(id)
	m, err := mid.Get(r.Context())
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	if err := m.SetChapters(r.Context(), chapters...); err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func getMemberVcard(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	mid := model.MemberID(id)
	member, err := mid.Get(r.Context())
	if err != nil {
		sendError(w, err, http.StatusNotFound)
		return
	}

	// Override global JSON header for vCard
	slog.Info("someone actually used the vard option!", "member", mid)
	w.Header().Set("Content-Type", "text/vcard")
	if err := member.WriteVCard(w); err != nil {
		slog.Error(err.Error())
		return
	}
}
