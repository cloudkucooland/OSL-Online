package model

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type GivingRecord struct {
	EntryID     int
	ID          MemberID
	Amount      float64
	Check       int
	Transaction string
	Description string
	Date        time.Time
}

func (n *GivingRecord) Store(ctx context.Context) error {
	zt, _ := time.Parse(timeformat, zerotime)
	if n.Date == zt {
		n.Date = time.Now()
	}
	_, err := db.ExecContext(ctx, "INSERT INTO `giving` (`entryID`, `id`, `amount`, `check`, `transaction`, `description`, `date`) VALUES (0,?,?,?,?,?,?)", n.ID, n.Amount, n.Check, makeNullString(n.Transaction), n.Description, makeNullTime(n.Date))
	if err != nil {
		slog.Error("database error in GivingRecord.Store", "err", err, "id", n.ID)
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (id MemberID) GivingRecords(ctx context.Context) ([]*GivingRecord, error) {
	gr := make([]*GivingRecord, 0)

	rows, err := db.QueryContext(ctx, "SELECT `entryID`, `amount`, `check`, `transaction`, `description`, `date` FROM `giving` WHERE `id` = ? ORDER BY `date`", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return gr, nil
		}
		slog.Error("database error in GivingRecords", "err", err, "id", id)
		return gr, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var g GivingRecord
		err := rows.Scan(&g.EntryID, &g.Amount, &g.Check, &g.Transaction, &g.Description, &g.Date)
		if err != nil {
			slog.Error("failed to scan row in GivingRecords", "err", err, "id", id)
			continue
		}
		gr = append(gr, &g)
	}
	return gr, nil
}
