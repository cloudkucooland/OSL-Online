package model

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func GetAuthData(id string) (string, int, error) {
	var pwhash string
	var level int
	err := db.QueryRow("SELECT pwhash, level FROM auth WHERE ID = ?", id).Scan(&pwhash, &level)
	if err != nil && err == sql.ErrNoRows {
		err = fmt.Errorf("user not found")
		slog.Error(err.Error(), "id", id)
		return "", 0, err
	}
	if err != nil {
		slog.Error(err.Error())
		return "", 0, err
	}
	return pwhash, level, nil
}
