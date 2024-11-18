package model

import (
	"time"
)

type NoteID uint

type Note struct {
	ID     NoteID
	member MemberID
	date   time.Time
	note   string
}

func (m MemberID) GetNotes() ([]*Note, error) {
	notes := make([]*Note, 0)

	// TODO XXX

	return notes, nil
}

func (n *Note) Store() error {
	// TODO XXX

	return nil
}

func (n NoteID) Delete() error {
	// TODO XXX

	return nil
}
