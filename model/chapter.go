package model

import (
	"database/sql"
	"log/slog"
)

type Chapter struct {
	ID    int
	Name  string
	Prior int
}

func (c *Chapter) Store() error {
	_, err := db.Exec("REPLACE INTO `chapters` (`id`, `name`, `prior` VALUES (?,?,?)", c.ID, c.Name, c.Prior)
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

func Chapters() ([]Chapter, error) {
	ch := make([]Chapter, 0)

	rows, err := db.Query("SELECT `id`, `name`, `prior` FROM `chapters` ORDER BY `name`")
	if err != nil && err == sql.ErrNoRows {
		return ch, nil
	}
	if err != nil {
		slog.Error(err.Error())
		return ch, err
	}

	var c Chapter
	for rows.Next() {
		err := rows.Scan(&c.ID, &c.Name, &c.Prior)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		ch = append(ch, c)
	}
	return ch, nil
}

func (c *Chapter) Members() ([]Member, error) {
	members := make([]Member, 0)

	rows, err := db.Query("SELECT m.ID, m.MemberStatus, m.FirstName, m.MiddleName, m.LastName, m.PreferredName, m.Title, m.LifevowName, m.Suffix, m.Address, m.AddressLine2, m.City, m.State, m.Country, m.PostalCode, m.PrimaryPhone, m.PrimaryEmail, m.Leadership, m.ListInDirectory, m.ListAddress, m.ListPrimaryPhone, m.ListPrimaryEmail FROM member=m, chaptermembers=x WHERE x.chapter = ? AND m.ID = x.member AND m.MemberStatus != 'Removed'  ORDER BY m.LastName", c.ID)
	if err != nil && err == sql.ErrNoRows {
		return members, nil
	}
	if err != nil {
		slog.Error(err.Error())
		return members, err
	}

	for rows.Next() {
		var n Member
		var lvn, pn, middle, suffix, line2, pphone, pcode, country, city, state, title, pemail, addr sql.NullString

		err := rows.Scan(&n.ID, &n.MemberStatus, &n.FirstName, &middle, &n.LastName, &pn, &title, &lvn, &suffix, &addr, &line2, &city, &state, &country, &pcode, &pphone, &pemail, &n.Leadership, &n.ListInDirectory, &n.ListAddress, &n.ListPrimaryPhone, &n.ListPrimaryEmail)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		if !n.ListInDirectory {
			continue
		}

		if lvn.Valid {
			n.LifevowName = lvn.String
		}
		if pn.Valid {
			n.PreferredName = pn.String
		}
		if middle.Valid {
			n.MiddleName = middle.String
		}
		if suffix.Valid {
			n.Suffix = suffix.String
		}
		if line2.Valid {
			n.AddressLine2 = line2.String
		}
		if pphone.Valid {
			n.PrimaryPhone = pphone.String
		}
		if pcode.Valid {
			n.PostalCode = pcode.String
		}
		if country.Valid {
			n.Country = country.String
		}
		if city.Valid {
			n.City = city.String
		}
		if state.Valid {
			n.State = state.String
		}
		if title.Valid {
			n.Title = title.String
		}
		if pemail.Valid {
			n.PrimaryEmail = pemail.String
		}
		if addr.Valid {
			n.Address = addr.String
		}

		if !n.ListAddress {
			n.Address = ""
			n.AddressLine2 = ""
			n.City = ""
			n.State = ""
			n.PostalCode = ""
		}

		if !n.ListPrimaryEmail {
			n.PrimaryEmail = ""
		}

		if !n.ListPrimaryPhone {
			n.PrimaryPhone = ""
		}

		members = append(members, n)
	}
	return members, nil
}
