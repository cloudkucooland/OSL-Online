package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// Member is the format sent to the UI
/* type Member struct {
	ID                 int
	MemberStatus       string
	FirstName          string
	MiddleName         string
	LastName           string
	PreferredName      string
	Title              string
	LifevowName        string
	Suffix             string
	Address            string
	AddressLine2       string
	City               string
	State              string
	Country            string
	PostalCode         string
	PrimaryPhone       string
	SecondaryPhone     string
	PrimaryEmail       string
	SecondaryEmail     string
	BirthDate          time.Time
	DateRecordCreated  time.Time
	Chapter            string
	DateFirstVows      time.Time
	DateReaffirmation  time.Time
	DateRemoved        time.Time
	DateDeceased       time.Time
	DateNovitiate      time.Time
	DateLifeVows       time.Time
	Status             string
	Leadership         string
	HowJoined          string
	HowRemoved         string
	ListInDirectory    bool
	ListAddress        bool
	ListPrimaryPhone   bool
	ListSecondaryPhone bool
	ListPrimaryEmail   bool
	ListSecondaryEmail bool
	Doxology           string
	Newsletter         string
	Communication      string
	Occupation         string
	Employer           string
	Denomination       string
} */

func SetMeField(id int, field string, value string) error {

	slog.Info("self-updating", "id", id, "field", field, "value", value)
	// XXX disable for now, need to restrict which fields people can change
	return nil

	if field == "id" {
		err := fmt.Errorf("cannot change ID")
		slog.Error(err.Error())
		return err
	}
	if strings.ContainsAny(field, "`;%") {
		err := fmt.Errorf("sql injection attempt [%s]", field)
		slog.Error(err.Error())
		return err
	}
	q := fmt.Sprintf("UPDATE `member` SET `%s` = ? WHERE `id` = ?", field)

	switch field {
	case "ListInDirectory", "ListAddress", "ListPrimaryPhone", "ListSecondaryPhone", "ListPrimaryEmail", "ListSecondaryEmail":
		var nb sql.NullBool
		nb.Valid = true
		nb.Bool = value == "true"
		if _, err := db.Exec(q, nb, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	case "BirthDate", "DateRecordCreated", "DateFirstVows", "DateReaffirmation", "DateRemoved", "DateDeceased", "DateNovitiate", "DateLifeVows":
		value = strings.TrimSpace(value)
		if value == "" {
			value = "0001-01-01"
		}
		t, err := time.Parse(format, value)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		if _, err := db.Exec(q, t, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	default:
		var ns sql.NullString
		value = strings.TrimSpace(value)
		if value == "" {
			ns.Valid = false
			ns.String = ""
		} else {
			ns.Valid = true
			ns.String = value
		}
		if _, err := db.Exec(q, ns, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	// XXX get ID
	if _, err := db.Exec("INSERT INTO auditlog VALUES (?, ?, ?, ?, CURRENT_DATE())", 0, id, field, value); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
