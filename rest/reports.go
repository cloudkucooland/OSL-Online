package rest

import (
	"encoding/csv"
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

	headers(w, r)
	switch report {
	case "avery":
		reportAvery(w, r, ps)
	case "email":
		reportEmail(w, r, ps)
	case "expired":
		reportExpired(w, r, ps)
	case "life":
		reportLife(w, r, ps)
	case "notrenewed":
		reportNotRenewed(w, r, ps)
	default:
		reportLife(w, r, ps)
	}
}

func reportNotRenewed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	out := [][]string{
		{"DateReaffirmation", "FirstName", "LastName", "PreferredName", "Title", "Address", "AddressLine2", "City", "State", "Country", "PostalCode", "PrimaryEmail"},
	}

	m, err := model.ReportNotRenewed()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	for _, n := range m {
		y, m, d := n.DateReaffirmation.Date()
		s := fmt.Sprintf("%04d-%02d-%02d", y, m, d)
		member := []string{s, n.FirstName, n.LastName, n.PreferredName, n.Title, n.Address, n.AddressLine2, n.City, n.State, n.Country, n.PostalCode, n.PrimaryEmail}
		out = append(out, member)
	}

	w.Header().Set("Content-Type", "text/csv")
	report := csv.NewWriter(w)
	report.WriteAll(out)

	if err := report.Error(); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func reportExpired(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	out := [][]string{
		{"DateReaffirmation", "FirstName", "LastName", "PreferredName", "Title", "Address", "AddressLine2", "City", "State", "Country", "PostalCode", "PrimaryEmail"},
	}

	m, err := model.ReportExpired()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	for _, n := range m {
		y, m, d := n.DateReaffirmation.Date()
		s := fmt.Sprintf("%04d-%02d-%02d", y, m, d)
		member := []string{s, n.FirstName, n.LastName, n.PreferredName, n.Title, n.Address, n.AddressLine2, n.City, n.State, n.Country, n.PostalCode, n.PrimaryEmail}
		out = append(out, member)
	}

	w.Header().Set("Content-Type", "text/csv")
	report := csv.NewWriter(w)
	report.WriteAll(out)

	if err := report.Error(); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func reportEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	out := [][]string{
		{"OSLName", "MemberStatus", "FirstName", "LastName", "PreferredName", "Title", "LifevowName", "Suffix", "PrimaryEmail", "SecondaryEmail", "ListPrimaryEmail", "ListSecondaryEmail", "Doxology", "Newsletter", "Communication"},
	}

	m, err := model.ActiveMemberIDs()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	for _, id := range m {
		n, err := id.Get(true)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, jsonError(err), http.StatusInternalServerError)
			return
		}
		oslName := n.OSLName()
		member := []string{oslName, n.MemberStatus, n.FirstName, n.LastName, n.PreferredName, n.Title, n.LifevowName, n.Suffix, n.PrimaryEmail, n.SecondaryEmail, fmt.Sprintf("%t", n.ListPrimaryEmail), fmt.Sprintf("%t", n.ListSecondaryEmail), n.Doxology, n.Newsletter, n.Communication}
		out = append(out, member)
	}

	w.Header().Set("Content-Type", "text/csv")
	report := csv.NewWriter(w)
	report.WriteAll(out)

	if err := report.Error(); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func reportAnnual(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	out := [][]string{
		{"OSLName", "OSLShortName", "FirstName", "LastName", "PreferredName", "Title", "Suffix", "FormattedAddress", "Doxology", "Newsletter", "Communication"},
	}

	m, err := model.ReportAnnual()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	for _, n := range m {
		addr, _ := n.FormatAddress()
		member := []string{n.OSLName(), n.OSLShortName(), n.FirstName, n.LastName, n.PreferredName, n.Title, n.Suffix, addr, n.Doxology, n.Newsletter, n.Communication}
		out = append(out, member)
	}

	w.Header().Set("Content-Type", "text/csv")
	report := csv.NewWriter(w)
	report.WriteAll(out)

	if err := report.Error(); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func reportLife(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	out := [][]string{
		{"OSLName", "OSLShortName", "FirstName", "LastName", "PreferredName", "Title", "Suffix", "LifevowName", "FormattedAddress", "Doxology", "Newsletter", "Communication"},
	}

	m, err := model.ReportLife()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	for _, n := range m {
		addr, _ := n.FormatAddress()
		member := []string{n.OSLName(), n.OSLShortName(), n.FirstName, n.LastName, n.PreferredName, n.Title, n.Suffix, n.LifevowName, addr, n.Doxology, n.Newsletter, n.Communication}
		out = append(out, member)
	}

	w.Header().Set("Content-Type", "text/csv")
	report := csv.NewWriter(w)
	report.WriteAll(out)

	if err := report.Error(); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
