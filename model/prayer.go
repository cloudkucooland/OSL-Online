package model

import (
	"context"
	"log/slog"
	"time"
)

type PrayerID int

type Prayer struct {
	PrayerID  PrayerID
	MemberID  MemberID
	OSLName   string
	Content   string
	Anonymous bool
	Date      time.Time
}

// GetPrayers fetches the list.
// If public is true, it masks names for anonymous entries.
func GetPrayers(ctx context.Context, memberID *MemberID, public bool) ([]Prayer, error) {
	prayers := make([]Prayer, 0)

	query := `SELECT p.id, p.member, p.content, p.anonymous, p.date, m.FirstName, m.LastName, m.Title, m.LifevowName, m.Suffix, m.PreferredName FROM prayers p JOIN member m ON p.member = m.id`

	args := []any{}
	if memberID != nil {
		query += ` WHERE p.member = ?`
		args = append(args, *memberID)
	}
	query += ` ORDER BY p.date DESC`
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p Prayer
		var n memberNulls
		var anon int

		err := rows.Scan(
			&p.PrayerID, &p.MemberID, &p.Content, &anon, &p.Date,
			&n.FirstName, &n.LastName, &n.Title, &n.LifevowName, &n.Suffix, &n.PreferredName,
		)
		if err != nil {
			return nil, err
		}

		p.Anonymous = (anon == 1)

		if public && p.Anonymous {
			p.OSLName = "A sibling or friend"
		} else {
			m := n.toMember()
			p.OSLName = m.OSLName()
		}

		prayers = append(prayers, p)
	}
	return prayers, nil
}

func (p *Prayer) Insert(ctx context.Context) error {
	requestingMemberID, err := IDFromContext(ctx)
	if err != nil {
		return err
	}

	slog.Info("adding prayer", "by", requestingMemberID, "p", p.Content)

	// date can be NULL and will default to now
	query := "INSERT INTO prayers (member, content, anonymous) VALUES (?, ?, ?)"
	res, err := db.ExecContext(ctx, query, requestingMemberID, p.Content, p.Anonymous)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	p.PrayerID = PrayerID(id)
	return nil
}

func DeletePrayer(ctx context.Context, id PrayerID) error {
	level := LevelFromContext(ctx)
	requestingMemberID, err := IDFromContext(ctx)
	if err != nil {
		return err
	}

	// Security check: Only owner or admin can delete
	query := "DELETE FROM prayers WHERE id = ?"
	if level < AuthLevelManager {
		query += " AND member = ?"
		_, err := db.ExecContext(ctx, query, id, requestingMemberID)
		return err
	}

	slog.Info("deleting prayer", "prayerID", id, "by", requestingMemberID)

	// 0 results deleted is not an error, this can be silent
	_, err = db.ExecContext(ctx, query, id)
	return err
}
