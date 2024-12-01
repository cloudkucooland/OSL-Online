package model

import (
	"database/sql"
	"log/slog"
)

const friendzoneMinDays = 365 * 5 // 5 years

// Friendzone moves all members who have not reaffirmed vows in the past friendzoneMinDays and converts
// them to friends
func Friendzone() error {
	rows, err := db.Query("SELECT ID FROM member WHERE MemberStatus = 'Annual Vows' AND DateReaffirmation < DATE_SUB(CURRENT_DATE(), INTERVAL ? DAY) ORDER BY LastName", friendzoneMinDays)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id MemberID
		if err := rows.Scan(&id); err != nil {
			slog.Error(err.Error())
			continue
		}

		f, err := id.Get()
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		slog.Info("friendzoning", "name", f.OSLName(), "last reaffirmation", f.DateReaffirmation)

		if err := id.makeFriend(); err != nil {
			slog.Error(err.Error())
			continue
		}
	}
	return nil
}

// the logic for moving to friend is already in SetMemberField, use that
func (id MemberID) makeFriend() error {
	return id.SetMemberField("MemberStatus", "Friend", MemberID(0))
}
