package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

const contentType = "Content-Type"
const csvMime = "text/csv;charset=utf-8;"
const pdfMime = "application/pdf;"

// const contentDisposition = "Content-Disposition"
// const csvDisposition = "attachment"
// const pdfDisposition = `attachment;filename="avery.pdf";`

func reports(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	report := ps.ByName("report")
	if report == "" {
		err := fmt.Errorf("report request not set")
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}
	slog.Info("report", "requested", report, "requester", getUser(r))

	var err error
	headers(w, r)
	switch report {
	case "avery":
		w.Header().Set(contentType, pdfMime)
		// w.Header().Set(contentDisposition, pdfDisposition)
		err = model.ReportAvery(w)
	case "annual":
		w.Header().Set(contentType, csvMime)
		err = model.ReportAnnual(w)
	case "email":
		w.Header().Set(contentType, csvMime)
		err = model.ReportAllEmail(w)
	case "expired":
		w.Header().Set(contentType, csvMime)
		err = model.ReportExpired(w)
	case "life":
		w.Header().Set(contentType, csvMime)
		err = model.ReportLife(w)
	case "doxprint":
		w.Header().Set(contentType, csvMime)
		err = model.DoxologyPrinted(w)
	case "doxemail":
		w.Header().Set(contentType, csvMime)
		err = model.DoxologyEmailed(w)
	case "fontemail":
		w.Header().Set(contentType, csvMime)
		err = model.ReportFontEmailed(w)
	case "allsubscribers":
		w.Header().Set(contentType, csvMime)
		err = model.ReportAllSubscribers(w)
	case "barb":
		w.Header().Set(contentType, csvMime)
		err = model.ReportBarb(w)
	default:
		w.Header().Set(contentType, csvMime)
		err = model.ReportBarb(w)
	}
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}
