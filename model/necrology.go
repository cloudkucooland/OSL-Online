package model

import (
	"database/sql"
	"log/slog"
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
