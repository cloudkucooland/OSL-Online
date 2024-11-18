package model

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Boostport/address"
)

type Locality struct {
	CountryCode  string
	LocalityCode string
	JointCode    string
	Locality     string
}

func Localities() ([]*Locality, error) {
	localities := make([]*Locality, 0)

	rows, err := db.Query("SELECT DISTINCT Country, State FROM member WHERE MemberStatus != 'Removed' AND MemberStatus != 'Deceased' AND ListInDirectory = true ORDER BY Country, State")
	if err != nil && err == sql.ErrNoRows {
		return localities, nil
	}
	if err != nil {
		slog.Error(err.Error())
		return localities, err
	}
	defer rows.Close()

	for rows.Next() {
		var l Locality
		var cc, lc sql.NullString

		if err := rows.Scan(&cc, &lc); err != nil {
			slog.Error(err.Error())
			continue
		}

		if !cc.Valid {
			continue
		}

		l.CountryCode = cc.String
		l.JointCode = cc.String

		if lc.Valid {
			l.JointCode = fmt.Sprintf("%s-%s", cc.String, lc.String)
			l.LocalityCode = lc.String
			l.Locality = lc.String

			country := address.GetCountry(l.CountryCode)

			// get friendly name for the locality
			if aa, ok := country.AdministrativeAreas["en"]; ok {
				for _, k := range aa {
					if k.ID == lc.String {
						l.Locality = k.Name
						break
					}
				}
			}
		}

		localities = append(localities, &l)
	}
	return localities, nil
}

func LocalityMembers(country string, locality string) ([]*Member, error) {
	if country == "SG" {
		return localityMembersSG()
	}

	members := make([]*Member, 0)

	var lc sql.NullString
	if locality != "" {
		lc.Valid = true
		lc.String = locality
	}

	rows, err := db.Query("SELECT ID FROM member WHERE Country = ? AND State = ? AND MemberStatus != 'Deceased' AND MemberStatus != 'Removed' AND ListInDirectory = true ORDER BY LastName, FirstName", country, lc)
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
			slog.Error(err.Error())
			continue
		}
		m.CleanUnlisted()
		members = append(members, m)
	}
	return members, nil
}

func localityMembersSG() ([]*Member, error) {
	members := make([]*Member, 0)

	rows, err := db.Query("SELECT ID FROM member WHERE Country = 'SG' AND MemberStatus != 'Deceased' AND MemberStatus != 'Removed' AND ListInDirectory = true ORDER BY LastName, FirstName")
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
			slog.Error(err.Error())
			continue
		}
		m.CleanUnlisted()
		members = append(members, m)
	}
	return members, nil
}
