package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudkucooland/OSL-Online/email"
	"github.com/cloudkucooland/OSL-Online/model"
)

func getMemberGiving(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}
	mid := model.MemberID(id)
	gr, err := mid.GivingRecords(r.Context())
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sendJSON(w, gr)
}

func postMemberGiving(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := r.ParseMultipartForm(1024 * 2); err != nil {
		slog.Warn(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	s := r.PostFormValue("id")
	if s == "" {
		sendError(w, fmt.Errorf("id not set"), http.StatusNotAcceptable)
		return
	}
	id, err := parseIDFromString(s)
	if err != nil {
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	s = r.PostFormValue("amount")
	if s == "" {
		sendError(w, fmt.Errorf("amount not set"), http.StatusNotAcceptable)
		return
	}
	amount, err := strconv.ParseFloat(s, 64)
	if err != nil {
		sendError(w, err, http.StatusNotAcceptable)
		return
	}

	s = r.PostFormValue("check")
	check := 0
	if s != "" {
		check, err = parseIDFromString(s)
		if err != nil {
			sendError(w, err, http.StatusNotAcceptable)
			return
		}
	}

	transaction := r.PostFormValue("transaction")
	description := r.PostFormValue("description")

	d := r.PostFormValue("date")
	date, err := time.Parse("2006-01-02", d)
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	gr := model.GivingRecord{
		EntryID:     0,
		ID:          model.MemberID(id),
		Amount:      amount,
		Check:       check,
		Transaction: transaction,
		Description: description,
		Date:        date,
	}

	if err := gr.Store(r.Context()); err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	if err := email.SendGiving(r.Context(), model.MemberID(id), fmt.Sprintf("%.2f", gr.Amount), gr.Description); err != nil {
		slog.Error(err.Error())
	}

	fmt.Fprint(w, jsonStatusOK)
}
