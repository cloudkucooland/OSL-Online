package model

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

func Necrology(ctx context.Context) ([]*Member, error) {
	members := make([]*Member, 0)

	rows, err := db.QueryContext(ctx, "SELECT ID FROM member WHERE MemberStatus = 'Deceased' AND DateDeceased IS NOT NULL ORDER BY LastName, FirstName")
	if err != nil {
		if err == sql.ErrNoRows {
			return members, nil
		}
		slog.Error("database error in Necrology", "err", err)
		return members, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id MemberID

		if err := rows.Scan(&id); err != nil {
			slog.Error("failed to scan row in Necrology", "err", err)
			continue
		}
		m, err := id.Get(ctx)
		if err != nil {
			continue
		}
		members = append(members, m)
	}
	return members, nil
}

type Commemoration struct {
	OSLName  string
	Locality string
	Country  string
	Year     int
}

func Commemorations(ctx context.Context, month time.Month, day int) ([]Commemoration, error) {
	commemorations := make([]Commemoration, 0)

	qq := fmt.Sprintf("%%-%02d-%02d", month, day)

	rows, err := db.QueryContext(ctx, "SELECT ID FROM member WHERE DateDeceased LIKE ? AND MemberStatus = 'Deceased' ORDER BY LastName, FirstName", qq)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Debug("no commemorations for this day", "month", month, "day", day)
			return commemorations, nil
		}
		slog.Error("database error in Commemorations", "err", err, "month", month, "day", day)
		return commemorations, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id MemberID
		err := rows.Scan(&id)
		if err != nil {
			slog.Error("failed to scan row in Commemorations", "err", err)
			continue
		}

		m, err := id.Get(ctx)
		if err != nil {
			continue
		}

		c := Commemoration{
			OSLName:  m.OSLName(),
			Locality: m.State,
			Country:  m.Country,
			Year:     m.DateDeceased.Year(),
		}

		commemorations = append(commemorations, c)
	}
	return commemorations, nil
}
