package model

import (
	"database/sql"
	"log/slog"
)

type ChapterID int

type Chapter struct {
	ID    ChapterID
	Name  string
	Prior MemberID
}

// don't use this... you'll lose all member data, rewrite to be an INSERT/ON DUPLICATE
func (c *Chapter) store() error {
	_, err := db.Exec("REPLACE INTO `chapters` (`id`, `name`, `prior` VALUES (?,?,?)", c.ID, c.Name, c.Prior)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (c *Chapter) Update() error {
	_, err := db.Exec("UPDATE `chapters` SET `name` = ?, `prior` = ?  WHERE `id` = ? `name`, `prior`", c.Name, c.Prior, c.ID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (c *Chapter) Remove() error {
	_, err := db.Exec("DELETE FROM `chapters` WHERE `ID` = ?", c.ID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func Chapters() ([]*Chapter, error) {
	ch := make([]*Chapter, 0)

	rows, err := db.Query("SELECT `id`, `name`, `prior` FROM `chapters` ORDER BY `name`")
	if err != nil && err == sql.ErrNoRows {
		return ch, nil
	}
	if err != nil {
		slog.Error(err.Error())
		return ch, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Chapter
		err := rows.Scan(&c.ID, &c.Name, &c.Prior)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		ch = append(ch, &c)
	}
	return ch, nil
}

func (c *Chapter) Members() ([]*Member, error) {
	members := make([]*Member, 0)

	rows, err := db.Query("SELECT m.ID FROM member=m, chaptermembers=x WHERE x.chapter = ? AND m.ID = x.member AND m.MemberStatus != 'Removed' AND m.MemberStatus != 'Deceased' ORDER BY m.LastName", c.ID)
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
		m, err := id.Get()
		if err != nil {
			continue
		}
		if !m.ListInDirectory {
			continue
		}
		m.CleanUnlisted()

		members = append(members, m)
	}
	return members, nil
}
