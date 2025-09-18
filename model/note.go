package model

import (
	"log/slog"
	"time"
)

type NoteID uint

type Note struct {
	ID     NoteID
	Member MemberID
	Date   time.Time
	Note   string
}

func (m MemberID) GetNotes() ([]*Note, error) {
	notes := make([]*Note, 0)

	rows, err :=  db.Query("SELECT ID, member, date, note FROM notes WHERE member = ?", m)
	if err != nil {
		slog.Error(err.Error())
		return notes, err
	}
	defer rows.Close()

	for rows.Next() {
		var n Note
		if err = rows.Scan(&n.ID, &n.Member, &n.Date, &n.Note); err != nil {
			slog.Error(err.Error())
			continue
		}
		notes = append(notes, &n)
	}
	return notes, nil
}

func (n *Note) Store() error {
	slog.Info("adding note", "member", n.Member, "value", n.Note)

	_, err := db.Query("INSERT INTO notes VALUES (0, ?, CURDATE(), ?)", n.Member, n.Note)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func (n NoteID) Delete() error {
	slog.Info("deleting note", "value", n)

	_, err := db.Query("DELETE FROM notes WHERE ID = ?", n)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
