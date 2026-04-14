package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"
)

func getMember(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("invalid id", "val", idStr, "err", err)
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	mid := model.MemberID(id)
	m, err := mid.Get(r.Context())
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	m.FormattedAddr, err = model.FormatAddress(m)
	if err != nil {
		slog.Error("address format failed", "err", err)
	}

	slog.Info("user loaded", "user", m.OSLName(), "requester", getUser(r))
	json.NewEncoder(w).Encode(m)
}

func getMemberChapters(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	mid := model.MemberID(id)
	chapters, err := mid.GetChapters()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(chapters)
}

func setMember(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	field := r.PostFormValue("field")
	value := r.PostFormValue("value")

	if field == "" {
		http.Error(w, jsonError(fmt.Errorf("field not set")), http.StatusNotAcceptable)
		return
	}

	if err := model.MemberID(id).SetMemberField(r.Context(), field, value); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func createMember(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	firstname := r.PostFormValue("firstname")
	lastname := r.PostFormValue("lastname")

	if firstname == "" || lastname == "" {
		http.Error(w, jsonError(fmt.Errorf("name components missing")), http.StatusNotAcceptable)
		return
	}

	id, err := model.Create(firstname, lastname)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `{"status":"ok", "id": %d}`, id)
}

func setMemberChapters(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, jsonError(err), http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(1024)
	cs := r.PostFormValue("chapters")
	var chapters []int

	if cs != "" {
		ss := strings.Split(cs, ",")
		for _, n := range ss {
			if c, err := strconv.Atoi(strings.TrimSpace(n)); err == nil {
				chapters = append(chapters, c)
			}
		}
	}

	mid := model.MemberID(id)
	m, err := mid.Get(r.Context())
	if err != nil {
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := m.SetChapters(r.Context(), chapters...); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func getMemberVcard(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	mid := model.MemberID(id)
	member, err := mid.Get(r.Context())
	if err != nil {
		http.Error(w, jsonError(err), http.StatusNotFound)
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
