package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
)

const contentType = "Content-Type"
const csvMime = "text/csv;charset=utf-8;"
const pdfMime = "application/pdf;"

func reports(w http.ResponseWriter, r *http.Request) {
	report := r.PathValue("report")
	if report == "" {
		err := fmt.Errorf("report request not set")
		slog.Error(err.Error())
		sendError(w, err, http.StatusNotAcceptable)
		return
	}
	slog.Info("report", "requested", report, "requester", getUser(r))

	var err error
	switch report {
	case "avery":
		w.Header().Set(contentType, pdfMime)
		err = model.ReportAvery(r.Context(), w)
	case "annual":
		w.Header().Set(contentType, csvMime)
		err = model.ReportAnnual(r.Context(), w)
	case "reaffirmation":
		w.Header().Set(contentType, csvMime)
		err = model.ReportReaffirmationFormMerge(r.Context(), w)
	case "email":
		w.Header().Set(contentType, csvMime)
		err = model.ReportAllEmail(r.Context(), w)
	case "expired":
		w.Header().Set(contentType, csvMime)
		err = model.ReportExpired(r.Context(), w)
	case "life":
		w.Header().Set(contentType, csvMime)
		err = model.ReportLife(r.Context(), w)
	case "lifecheckin":
		w.Header().Set(contentType, csvMime)
		err = model.ReportLifeCheckinFormMerge(r.Context(), w)
	case "doxprint":
		w.Header().Set(contentType, csvMime)
		err = model.DoxologyPrinted(r.Context(), w)
	case "allsubscribers":
		w.Header().Set(contentType, csvMime)
		err = model.ReportAllSubscribers(r.Context(), w)
	case "barb":
		w.Header().Set(contentType, csvMime)
		err = model.ReportBarb(r.Context(), w)
	default:
		w.Header().Set(contentType, csvMime)
		err = model.ReportBarb(r.Context(), w)
	}
	if err != nil {
		slog.Error(err.Error())
		sendError(w, err, http.StatusInternalServerError)
		return
	}
}
