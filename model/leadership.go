package model

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

func Leadership(ctx context.Context, category string) ([]*Member, error) {
	members := make([]*Member, 0)

	if category == "" {
		slog.Warn("Leadership category is empty, defaulting to member")
		category = "member"
	}

	rows, err := db.QueryContext(ctx, "SELECT ID FROM member WHERE Leadership = ? AND MemberStatus NOT IN ('Removed', 'Deceased') ORDER BY LastName", category)
	if err != nil {
		if err == sql.ErrNoRows {
			return members, nil
		}
		slog.Error("database error in Leadership", "err", err, "category", category)
		return members, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id MemberID
		if err := rows.Scan(&id); err != nil {
			slog.Error("failed to scan row in Leadership", "err", err, "category", category)
			continue
		}
		m, err := id.Get(ctx)
		if err != nil {
			// id.Get already logs
			continue
		}
		members = append(members, m)
	}
	return members, nil
}
