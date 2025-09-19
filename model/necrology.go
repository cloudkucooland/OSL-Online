package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

func Necrology() ([]*Member, error) {
	members := make([]*Member, 0)

	rows, err := db.Query("SELECT ID FROM member WHERE MemberStatus = 'Deceased' AND DateDeceased IS NOT NULL ORDER BY LastName, FirstName")
	if err != nil && err == sql.ErrNoRows {
		return members, nil
	}
	if err != nil {
		slog.Error(err.Error())
		return members, err
	}
	defer rows.Close()

	for rows.Next() {
		var id MemberID

		if err := rows.Scan(&id); err != nil {
			slog.Error(err.Error())
			continue
		}
		m, err := id.Get()
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
}

func Commemorations(month time.Month, day int) ([]Commemoration, error) {
	commemorations := make([]Commemoration, 0)

	qq := fmt.Sprintf("%%-%02d-%02d", month, day)

	slog.Info("commemorations qq", "qq", qq)

	rows, err := db.Query("SELECT ID FROM member WHERE DateDeceased LIKE ? AND MemberStatus = 'Deceased' ORDER BY LastName, FirstName", qq)
	if err != nil && err == sql.ErrNoRows {
		slog.Info("no commemorations for this day", "month", month, "day", day)
		return commemorations, nil
	}
	if err != nil {
		slog.Error(err.Error())
		return commemorations, err
	}
	defer rows.Close()

	for rows.Next() {
		var id MemberID
		err := rows.Scan(&id)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		m, err := id.Get()
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		c := Commemoration{
			OSLName:  m.OSLName(),
			Locality: m.State,
			Country:  m.Country,
		}

		commemorations = append(commemorations, c)
	}
	return commemorations, nil
}
