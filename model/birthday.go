package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

func SearchBirthday(month time.Month, day int) ([]*Member, error) {
	members := make([]*Member, 0)

	qq := fmt.Sprintf("%%-%02d-%02d", month, day)

	rows, err := db.Query("SELECT ID FROM member WHERE BirthDate LIKE ? AND MemberStatus NOT IN ('Removed', 'Deceased') ORDER BY LastName, FirstName", qq)
	if err != nil && err == sql.ErrNoRows {
		slog.Info("no birthdays for this day", "month", month, "day", day)
		return members, nil
	}
	if err != nil {
		slog.Error(err.Error())
		return members, err
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

		members = append(members, m)
	}
	return members, nil
}
