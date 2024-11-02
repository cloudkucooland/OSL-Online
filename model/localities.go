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

// why doesn't Singapore show up in this?
func Localities() ([]Locality, error) {
	localities := make([]Locality, 0)

	rows, err := db.Query("SELECT Country, State, CONCAT(Country, '-', State) AS CS FROM member WHERE MemberStatus != 'Removed' AND MemberStatus != 'Deceased' AND ListInDirectory = true GROUP BY (CS) ORDER BY CS")
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
		var cc, lc, jc sql.NullString

		if err := rows.Scan(&cc, &lc, &jc); err != nil {
			slog.Error(err.Error())
			continue
		}

		if !cc.Valid {
			continue
		}

		l.CountryCode = cc.String
		l.JointCode = cc.String

		if lc.Valid {
			country := address.GetCountry(l.CountryCode)
			if aa, ok := country.AdministrativeAreas["en"]; ok {
				for _, k := range aa {
					if k.ID == lc.String {
						l.LocalityCode = lc.String
						l.Locality = k.Name
						l.JointCode = fmt.Sprintf("%s-%s", cc.String, lc.String)
						break
					}
				}
			}
		}

		localities = append(localities, l)
	}
	return localities, nil
}

func LocalityMembers(country string, locality string) ([]*Member, error) {
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

		m, err := id.Get(false)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		members = append(members, m)
	}
	return members, nil
}
