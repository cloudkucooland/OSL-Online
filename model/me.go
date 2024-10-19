package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

func SetMeField(id int, field string, value string) error {
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
	default:
		err := fmt.Errorf("cannot edit that field")
		slog.Error(err.Error(), "id", id, "field", field, "value", value)
		return err
	}

	if _, err := db.Exec("INSERT INTO auditlog VALUES (?, ?, ?, ?, CURRENT_DATE())", id, id, field, value); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
