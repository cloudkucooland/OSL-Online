package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func reports(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	report := ps.ByName("report")
	if report == "" {
		err := fmt.Errorf("report request not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}
	slog.Info("report", "requested", report, "requester", getUser(r))

	headers(w, r)
	switch report {
	case "avery":
		reportAvery(w, r, ps)
	case "annual":
		reportAnnual(w, r, ps)
	case "email":
		reportEmail(w, r, ps)
	case "expired":
		reportExpired(w, r, ps)
	case "life":
		reportLife(w, r, ps)
	case "doxprint":
		reportDoxPrinted(w, r, ps)
	case "doxemail":
		reportDoxEmailed(w, r, ps)
	case "fontemail":
		reportFontEmailed(w, r, ps)
	default:
		reportLife(w, r, ps)
	}
}

func reportExpired(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/csv")
	if err := model.ReportExpired(w); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func reportEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/csv")
	if err := model.ReportAllEmail(w); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func reportAnnual(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/csv")
	if err := model.ReportAnnual(w); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func reportLife(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/csv")
	if err := model.ReportLife(w); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func reportAvery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/pdf")
	if err := model.ReportAvery(w); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

// convert older reports to this new way of doing things
func reportDoxPrinted(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/csv")
	if err := model.DoxologyPrinted(w); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func reportDoxEmailed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/csv")
	if err := model.DoxologyEmailed(w); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func reportFontEmailed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/csv")
	if err := model.ReportFontEmailed(w); err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}
