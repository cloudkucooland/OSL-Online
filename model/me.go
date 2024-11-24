package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

func SetMeField(id MemberID, field string, value string) error {
	slog.Info("self-updating", "id", id, "field", field, "value", value)

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
	// These are allowed
	case "ListInDirectory", "ListAddress", "ListPrimaryPhone", "ListSecondaryPhone", "ListPrimaryEmail", "ListSecondaryEmail":
		var nb sql.NullBool
		nb.Valid = true
		nb.Bool = value == "true"
		if _, err := db.Exec(q, nb, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	// These are allowed
	case "BirthDate":
		value = strings.TrimSpace(value)
		if value == "" {
			value = zerotime
		}
		t, err := time.Parse(timeformat, value)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		if _, err := db.Exec(q, t, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	// These are allowed
	case "PreferredName", "Title", "Occupation", "Employer", "Denomination":
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
	case "Communication": // users can do printed even if not donated this year
		value = strings.TrimSpace(value)
		if _, err := db.Exec(q, value, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	case "Newsletter", "Doxology": // only allow printed if donated this year
		value = strings.TrimSpace(value)
		if value == "mailed" && !id.allowPrinted() {
			err := fmt.Errorf("no donations in the past 12 months, cannot choose printed Doxology or Newsletter")
			return err
		}
		if _, err := db.Exec(q, value, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	default:
		err := fmt.Errorf("cannot edit that field")
		slog.Error(err.Error(), "id", id, "field", field, "value", value)
		return err
	}

	if err := id.ChangeLogStore(ChangeLogEntry{
		Changer: id,
		Field:   field,
		Value:   value,
	}); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func (id MemberID) allowPrinted() bool {
	m, err := id.Get()
	if err != nil {
		return false
	}

	if m.MemberStatus == "Removed" || m.MemberStatus == "Deceased" {
		return false
	}

	gr, err := id.GivingRecords()
	if err != nil {
		return false
	}
	found := false
	yearago := time.Now().AddDate(-1, 0, 0)
	for _, r := range gr {
		if r.Date.After(yearago) {
			found = true
			break
		}
	}
	return found
}
