package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func getMemberGiving(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	mid := model.MemberID(id)
	m, err := mid.Get(false)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	gr, err := m.GivingRecords()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(gr); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func postMemberGiving(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)
	if err := r.ParseMultipartForm(1024 * 64); err != nil {
		slog.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	s := r.PostFormValue("id")
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

	s = r.PostFormValue("amount")
	if s == "" {
		err := fmt.Errorf("amount not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}
	amount, err := strconv.ParseFloat(s, 64)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	s = r.PostFormValue("check")
	check := 0
	if s != "" {
		check, err = strconv.Atoi(s)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, jsonError(err), http.StatusNotAcceptable)
			return
		}
	}

	transaction := r.PostFormValue("transaction")
	description := r.PostFormValue("description")

	d := r.PostFormValue("date")
	date, err := time.Parse("2006-01-02", d)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
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

	if err := gr.Store(); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
