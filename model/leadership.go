package model

import (
	"database/sql"
	"log/slog"
)

func Leadership(category string) ([]Member, error) {
	members := make([]Member, 0)

	rows, err := db.Query("SELECT ID, MemberStatus, FirstName, MiddleName, LastName, PreferredName, Title, LifevowName, Suffix, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, PrimaryEmail, Leadership, ListInDirectory, ListAddress, ListPrimaryPhone, ListPrimaryEmail FROM member WHERE Leadership = ? AND MemberStatus != 'Removed' ORDER BY LastName", category)
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
