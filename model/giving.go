package model

import (
	"database/sql"
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

func (n GivingRecord) Store() error {
	_, err := db.Exec("INSERT INTO `giving` (`entryID`, `id`, `amount`, `check`, `transaction`, `description`, `date`) VALUES (0,?,?,?,?,?,curdate())", n.ID, n.Amount, n.Check, makeNullString(n.Transaction), n.Description)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (id MemberID) GivingRecords() ([]GivingRecord, error) {
	gr := make([]GivingRecord, 0)

	rows, err := db.Query("SELECT `entryID`, `amount`, `check`, `transaction`, `description`, `date` FROM `giving` WHERE `id` = ? ORDER BY `date`", id)
	if err != nil && err == sql.ErrNoRows {
		return gr, nil
	}
	if err != nil {
		slog.Error(err.Error())
		return gr, err
	}
	defer rows.Close()

	var g GivingRecord
	var d string
	for rows.Next() {
		err := rows.Scan(&g.EntryID, &g.Amount, &g.Check, &g.Transaction, &g.Description, &d)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		g.Date, _ = time.Parse(timeformat, d)

		gr = append(gr, g)
	}
	return gr, nil
}
