package model

import (
	"database/sql"
	"fmt"
	"log/slog"
)

// SearchResult is the format sent to the UI
type SearchResult struct {
	ID              int
	MemberStatus    string
	FirstName       string
	LastName        string
	PreferredName   string
	ListInDirectory bool
}

func Search(query string, unlisted bool) ([]SearchResult, error) {
	var res []SearchResult
	var n SearchResult

	qq := fmt.Sprintf("%%%s%%", query)

	var pn sql.NullString

	rows, err := db.Query("SELECT ID, MemberStatus, FirstName, LastName, PreferredName, ListInDirectory FROM member WHERE FirstName like ? OR LastName like ? OR PreferredName LIKE ? OR LifeVowName LIKE ? ORDER BY LastName, FirstName", qq, qq, qq, qq)
	if err != nil {
		slog.Error(err.Error())
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&n.ID, &n.MemberStatus, &n.FirstName, &n.LastName, &pn, &n.ListInDirectory); err != nil {
			slog.Error(err.Error())
			continue
		}
		if pn.Valid {
			n.PreferredName = pn.String
		} else {
			n.PreferredName = ""
		}

		if n.PreferredName == n.FirstName {
			n.PreferredName = ""
		}
		if unlisted || n.ListInDirectory {
			res = append(res, n)
		}
	}

	return res, nil
}

func SearchEmail(query string, unlisted bool) ([]SearchResult, error) {
	var res []SearchResult
	var n SearchResult

	var pn sql.NullString

	rows, err := db.Query("SELECT ID, MemberStatus, FirstName, LastName, PreferredName, ListInDirectory FROM member WHERE PrimaryEmail = ? OR SecondaryEmail = ?", query, query)
	if err != nil {
		slog.Error(err.Error())
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&n.ID, &n.MemberStatus, &n.FirstName, &n.LastName, &pn, &n.ListInDirectory); err != nil {
			slog.Error(err.Error())
			continue
		}
		if pn.Valid {
			n.PreferredName = pn.String
		} else {
			n.PreferredName = ""
		}

		if n.PreferredName == n.FirstName {
			n.PreferredName = ""
		}
		if unlisted || n.ListInDirectory {
			res = append(res, n)
		}
	}

	return res, nil
}
