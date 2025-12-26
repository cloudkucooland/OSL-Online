package model

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

func SetMeField(ctx context.Context, id MemberID, field string, value string) error {
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
		if _, err := db.ExecContext(ctx, q, nb, id); err != nil {
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
		if _, err := db.ExecContext(ctx, q, t, id); err != nil {
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
		if _, err := db.ExecContext(ctx, q, ns, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	case "Communication": // users can do printed even if not donated this year
		cp := communicationPref(strings.TrimSpace(value))
		switch cp {
		case ELECTRONIC, MAILED:
		default:
			cp = NONE
		}

		if _, err := db.ExecContext(ctx, q, cp, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	case "Newsletter":
		cp := communicationPref(strings.TrimSpace(value))
		if cp == MAILED && !id.allowPrinted(ctx) {
			err := fmt.Errorf("no donations in the past 12 months, cannot choose printed Newsletter")
			return err
		}

		switch cp {
		case ELECTRONIC, MAILED:
			if err := id.SubscribeFont(ctx); err != nil {
				slog.Error(err.Error())
				// continue
			}
		default:
			cp = NONE
			if err := id.UnsubscribeFont(ctx); err != nil {
				slog.Error(err.Error())
				// continue
			}
		}

		if _, err := db.ExecContext(ctx, q, cp, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	case "Doxology": // only allow printed if donated this year
		cp := communicationPref(strings.TrimSpace(value))
		if cp == MAILED && !id.allowPrinted(ctx) {
			err := fmt.Errorf("no donations in the past 12 months, cannot choose printed Doxology")
			return err
		}

		switch cp {
		case ELECTRONIC, MAILED:
			if err := id.SubscribeDoxology(ctx); err != nil {
				slog.Error(err.Error())
				// continue
			}
		default:
			cp = NONE
			if err := id.UnsubscribeDoxology(ctx); err != nil {
				slog.Error(err.Error())
				// continue
			}
		}

		if _, err := db.ExecContext(ctx, q, cp, id); err != nil {
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

func (id MemberID) allowPrinted(ctx context.Context) bool {
	m, err := id.Get(ctx)
	if err != nil {
		return false
	}

	if m.MemberStatus == REMOVED || m.MemberStatus == DECEASED {
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
