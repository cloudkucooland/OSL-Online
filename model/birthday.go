package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

func SearchBirthday(month time.Month, day int) ([]*Member, error) {
	members := make([]*Member, 0)

	qq := fmt.Sprintf("%%-%2d-%2d", month, day)

	rows, err := db.Query("SELECT ID FROM member WHERE BirthDate LIKE ? AND MemberStatus != 'Removed' AND MemberStatus != 'Deceased' ORDER BY LastName, FirstName", qq)
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
		err := rows.Scan(&id)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		m, err := id.Get(true)
		if err != nil {
			slog.Error(err.Error())
			continue
			// return members, err
		}

		members = append(members, m)
	}
	return members, nil
}
