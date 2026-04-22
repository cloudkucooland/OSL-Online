package model

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type ChapterID int

type Chapter struct {
	ID    ChapterID
	Name  string
	Prior MemberID
	Email string
}

// don't use this... you'll lose all member data, rewrite to be an INSERT/ON DUPLICATE
/*
func (c *Chapter) store() error {
	_, err := db.Exec("REPLACE INTO `chapters` (`id`, `name`, `prior` VALUES (?,?,?)", c.ID, c.Name, c.Prior)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
} */

func (c *Chapter) Update(ctx context.Context) error {
	var e sql.NullString

	_, err := db.ExecContext(ctx, "UPDATE `chapters` SET `name` = ?, `prior` = ?, `email` = ? WHERE `id` = ?", c.Name, c.Prior, c.ID, e)
	if err != nil {
		slog.Error("database error in Chapter.Update", "err", err, "id", c.ID)
		return fmt.Errorf("database error: %w", err)
	}
	if e.Valid {
		c.Email = e.String
	}

	return nil
}

func (c *Chapter) Remove(ctx context.Context) error {
	_, err := db.ExecContext(ctx, "DELETE FROM `chapters` WHERE `ID` = ?", c.ID)
	if err != nil {
		slog.Error("database error in Chapter.Remove", "err", err, "id", c.ID)
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func Chapters(ctx context.Context) ([]*Chapter, error) {
	ch := make([]*Chapter, 0)

	rows, err := db.QueryContext(ctx, "SELECT `id`, `name`, `prior`, `email` FROM `chapters` ORDER BY `name`")
	if err != nil {
		if err == sql.ErrNoRows {
			return ch, nil
		}
		slog.Error("database error in Chapters", "err", err)
		return ch, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var e sql.NullString

	for rows.Next() {
		var c Chapter
		err := rows.Scan(&c.ID, &c.Name, &c.Prior, &e)
		if err != nil {
			slog.Error("failed to scan row in Chapters", "err", err)
			continue
		}
		if e.Valid {
			c.Email = e.String
		}

		ch = append(ch, &c)
	}
	return ch, nil
}

func (c *Chapter) Members(ctx context.Context) ([]*Member, error) {
	members := make([]*Member, 0)

	rows, err := db.QueryContext(ctx, "SELECT m.ID FROM member=m, chaptermembers=x WHERE x.chapter = ? AND m.ID = x.member AND m.MemberStatus NOT IN ('Removed', 'Deceased') ORDER BY m.LastName", c.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return members, nil
		}
		slog.Error("database error in Chapter.Members", "err", err, "id", c.ID)
		return members, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id MemberID

		if err := rows.Scan(&id); err != nil {
			slog.Error("failed to scan row in Chapter.Members", "err", err, "id", c.ID)
			continue
		}
		m, err := id.Get(ctx)
		if err != nil {
			continue
		}
		if !m.ListInDirectory {
			continue
		}

		members = append(members, m)
	}
	return members, nil
}

func (id ChapterID) Load(ctx context.Context) (*Chapter, error) {
	var c Chapter
	err := db.QueryRowContext(ctx, "SELECT `id`, `name`, `prior`, `email` FROM `chapters` WHERE `id` = ?", id).Scan(&c.ID, &c.Name, &c.Prior, &c.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Warn("chapter not found", "id", id)
			return nil, fmt.Errorf("chapter %d not found", id)
		}
		slog.Error("database error in ChapterID.Load", "err", err, "id", id)
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &c, nil
}
