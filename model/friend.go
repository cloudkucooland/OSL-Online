package model

import (
	"context"
	"database/sql"
	"log/slog"
)

const friendzoneMinDays = 365 * 3 // 3 years

// Friendzone moves all members who have not reaffirmed vows in the past friendzoneMinDays and converts
// them to friends
func Friendzone(ctx context.Context) error {
	rows, err := db.QueryContext(ctx, "SELECT ID FROM member WHERE MemberStatus = 'Annual Vows' AND DateReaffirmation < DATE_SUB(CURRENT_DATE(), INTERVAL ? DAY) ORDER BY LastName", friendzoneMinDays)
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

		f, err := id.Get(ctx)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		slog.Info("friendzoning", "name", f.OSLName(), "last reaffirmation", f.DateReaffirmation)

		if err := id.makeFriend(ctx); err != nil {
			slog.Error(err.Error())
			continue
		}
	}
	return nil
}

// the logic for moving to friend is already in SetMemberField, use that
func (id MemberID) makeFriend(ctx context.Context) error {
	return id.SetMemberField(ctx, "MemberStatus", "Friend", MemberID(0))
}
