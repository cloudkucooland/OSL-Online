package model

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
)

type ChangeLogEntry struct {
	ID      MemberID
	Changer MemberID
	Field   string
	Value   string
	Date    time.Time
}

func (m MemberID) Changelog(ctx context.Context) ([]*ChangeLogEntry, error) {
	cr := make([]*ChangeLogEntry, 0)

	rows, err := db.QueryContext(ctx, "SELECT `changee`, `changer`, `field`, `value`, `date` FROM `auditlog` WHERE `changee` = ? ORDER BY `date`", m)
	if err != nil && err == sql.ErrNoRows {
		return cr, nil
	}
	if err != nil {
		slog.Error(err.Error())
		return cr, err
	}
	defer rows.Close()

	for rows.Next() {
		var c ChangeLogEntry
		err := rows.Scan(&c.ID, &c.Changer, &c.Field, &c.Value, &c.Date)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		cr = append(cr, &c)
	}
	return cr, nil
}

func (id MemberID) ChangeLogStore(ctx context.Context, c ChangeLogEntry) error {
	if _, err := db.ExecContext(ctx, "INSERT INTO auditlog VALUES (?, ?, ?, ?, CURRENT_DATE())", c.Changer, id, c.Field, c.Value); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
