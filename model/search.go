package model

import (
	"fmt"
	"log/slog"
)

// SearchResult is the format sent to the UI
type SearchResult struct {
	ID           int
	MemberStatus string
	FirstName    string
	LastName     string
}

// GetMember returns a populated Member struct, NULLs converted to ""
func Search(query string) ([]SearchResult, error) {
	var res []SearchResult
	var n SearchResult

	qq := fmt.Sprintf("%%%s%%", query)

	rows, err := db.Query("SELECT ID, MemberStatus, FirstName, LastName FROM member WHERE FirstName like ? OR LastName like ? ORDER BY LastName, FirstName", qq, qq)
	if err != nil {
		slog.Error(err.Error())
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&n.ID, &n.MemberStatus, &n.FirstName, &n.LastName); err != nil {
			slog.Error(err.Error())
			continue
		}
		res = append(res, n)
	}

	return res, nil
}
