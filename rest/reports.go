package rest

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/julienschmidt/httprouter"
)

func reportNotrenewed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	out := [][]string{
		{"DateReaffirmation", "FirstName", "LastName", "PreferredName", "Title", "Address", "AddressLine2", "City", "State", "Country", "PostalCode", "PrimaryEmail"},
	}
	headers(w, r)

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
	headers(w, r)

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
		{"MemberStatus", "FirstName", "LastName", "PreferredName", "Title", "LifevowName", "Suffix", "PrimaryEmail", "SecondaryEmail", "ListPrimaryEmail", "ListSecondaryEmail", "Doxology", "Newsletter", "Communication"},
	}
	headers(w, r)

	m, err := model.ReportEmail()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	for _, n := range m {
		member := []string{n.MemberStatus, n.FirstName, n.LastName, n.PreferredName, n.Title, n.LifevowName, n.Suffix, n.PrimaryEmail, n.SecondaryEmail, fmt.Sprintf("%t", n.ListPrimaryEmail), fmt.Sprintf("%t", n.ListSecondaryEmail), n.Doxology, n.Newsletter, n.Communication}
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
		{"PostalName", "FirstName", "LastName", "PreferredName", "Title", "Suffix", "Address", "AddressLine2", "City", "State", "Country", "PostalCode", "Doxology", "Newsletter", "Communication"},
	}
	headers(w, r)

	m, err := model.ReportAnnual()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	for _, n := range m {
		postalName := n.OSLName()
		member := []string{postalName, n.FirstName, n.LastName, n.PreferredName, n.Title, n.Suffix, n.Address, n.AddressLine2, n.City, n.State, n.Country, n.PostalCode, n.Doxology, n.Newsletter, n.Communication}
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
		{"PostalName", "FirstName", "LastName", "PreferredName", "Title", "Suffix", "LifevowName", "Address", "AddressLine2", "City", "State", "Country", "PostalCode", "Doxology", "Newsletter", "Communication"},
	}
	headers(w, r)

	m, err := model.ReportLife()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	for _, n := range m {
		postalName := n.OSLName()
		member := []string{postalName, n.FirstName, n.LastName, n.PreferredName, n.Title, n.Suffix, n.LifevowName, n.Address, n.AddressLine2, n.City, n.State, n.Country, n.PostalCode, n.Doxology, n.Newsletter, n.Communication}
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
