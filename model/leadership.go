package model

import (
	"context"
	"database/sql"
	"log/slog"
)

func Leadership(ctx context.Context, category string) ([]*Member, error) {
	members := make([]*Member, 0)

	rows, err := db.QueryContext(ctx, "SELECT ID FROM member WHERE Leadership = ? AND MemberStatus NOT IN ('Removed', 'Deceased') ORDER BY LastName", category)
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
		m, err := id.Get(ctx)
		if err != nil {
			continue
		}
		m.CleanUnlisted()
		members = append(members, m)
	}
	return members, nil
}
