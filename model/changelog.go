package model

import (
	"database/sql"
	"log/slog"
	"time"
)

type ChangeLogEntry struct {
	ID      int
	Changer int
	Field   string
	Value   string
	Date    time.Time
}

func (m *Member) Changelog() ([]ChangeLogEntry, error) {
	cr := make([]ChangeLogEntry, 0)

	rows, err := db.Query("SELECT `changee`, `changer`, `field`, `value`, `date` FROM `auditlog` WHERE `changee` = ? ORDER BY `date`", m.ID)
	if err != nil && err == sql.ErrNoRows {
		return cr, nil
	}
	if err != nil {
		slog.Error(err.Error())
		return cr, err
	}

	var c ChangeLogEntry
	var d string
	for rows.Next() {
		err := rows.Scan(&c.ID, &c.Changer, &c.Field, &c.Value, &d)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		c.Date, _ = time.Parse(format, d)

		cr = append(cr, c)
	}
	return cr, nil
}
