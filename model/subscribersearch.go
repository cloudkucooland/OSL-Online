package model

import (
	"database/sql"
	"fmt"
	"log/slog"
)

// SearchResult is the format sent to the UI
type SubSearchResult struct {
	ID   int
	Name string
	Attn string
}

// GetMember returns a populated Member struct, NULLs converted to ""
func SubscriberSearch(query string) ([]*SubSearchResult, error) {
	res := make([]*SubSearchResult, 0)
	var attn sql.NullString

	qq := fmt.Sprintf("%%%s%%", query)

	rows, err := db.Query("SELECT ID, Name, Attn FROM subscriber WHERE Name like ? OR Attn like ? ORDER BY Name, Attn", qq, qq)
	if err != nil {
		slog.Error(err.Error())
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var n SubSearchResult
		if err = rows.Scan(&n.ID, &n.Name, &attn); err != nil {
			slog.Error(err.Error())
			continue
		}
		n.Attn = ""
		if attn.Valid {
			n.Attn = attn.String
		}
		res = append(res, &n)
	}

	return res, nil
}
