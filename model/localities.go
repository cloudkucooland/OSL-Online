package model

import (
	"database/sql"
	"log/slog"
)

type Locality struct {
	Country  string
	Locality string
}

func Localities() ([]Locality, error) {
	localities := make([]Locality, 0)

	rows, err := db.Query("SELECT Country, State, CONCAT(Country, '-', State) AS CS FROM member WHERE State IS NOT NULL GROUP BY (CS) ORDER BY CS")
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
		var junk string
		err := rows.Scan(&l.Country, &l.Locality, &junk)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		localities = append(localities, l)
	}
	return localities, nil
}
