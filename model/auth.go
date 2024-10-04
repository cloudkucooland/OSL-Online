package model

import (
	"database/sql"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/cloudkucooland/OSL-Online/email"
	"github.com/sethvargo/go-password/password"
)

func GetAuthData(id string) (string, int, error) {
	var pwhash string
	var level int
	err := db.QueryRow("SELECT pwhash, level FROM auth WHERE user = ?", id).Scan(&pwhash, &level)
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

func SetAuthData(id string, pw string, level int) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	if err != nil {
		return err
	}

	_, err = db.Exec("REPLACE INTO auth VALUES (?,?,?)", id, bytes, level)
	if err != nil {
		return err
	}
	return nil
}

func Register(addr string) error {
	password, err := password.Generate(10, 3, 0, false, true)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if err := SetAuthData(addr, password, 0); err != nil {
		slog.Error(err.Error())
		return err
	}

	if email.SendRegister(addr, password); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func GetID(addr string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM member WHERE PrimaryEmail = ?", addr).Scan(&id)
	if err != nil && err == sql.ErrNoRows {
		err = fmt.Errorf("unknown email address")
		slog.Error(err.Error(), "addr", addr)
		return 0, err
	}
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}
	return id, nil
}
