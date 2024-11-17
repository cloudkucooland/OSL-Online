package model

import (
	"database/sql"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/sethvargo/go-password/password"
)

type Authname string

// String satisfies the stringer interface
func (u Authname) String() string {
	return string(u)
}

func (u Authname) getAuthData() (string, uint8, error) {
	var pwhash string
	var level uint8
	err := db.QueryRow("SELECT pwhash, level FROM auth WHERE user = ?", u).Scan(&pwhash, &level)
	if err != nil && err == sql.ErrNoRows {
		err = fmt.Errorf("user %s not found", u)
		slog.Error(err.Error(), "username", u)
		return "", 0, err
	}
	if err != nil {
		slog.Error(err.Error())
		return "", 0, err
	}
	return pwhash, level, nil
}

func (u Authname) Authenticate(password string) (uint8, error) {
	pwhash, level, err := u.getAuthData()
	if err != nil || pwhash == "" {
		err := fmt.Errorf("the email address %s has not yet been registered", u)
		slog.Error(err.Error())
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pwhash), []byte(password)); err != nil {
		slog.Error("login failed", "err", err)
		return 0, err
	}

	return level, nil
}

func (u Authname) SetAuthData(pw string, level int) error {
	slog.Info("updating password", "username", u)
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	if err != nil {
		return err
	}

	_, err = db.Exec("REPLACE INTO auth VALUES (?,?,?)", u, bytes, level)
	if err != nil {
		return err
	}
	return nil
}

func (u Authname) Register() (string, error) {
	slog.Info("registering user", "username", u)
	if _, err := u.GetID(); err != nil {
		return "", err
	}

	password, err := password.Generate(10, 3, 0, false, true)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	if err := u.SetAuthData(password, 0); err != nil {
		slog.Error(err.Error())
		return password, err
	}

	// caller must send email
	return password, nil
}

func (u Authname) GetID() (MemberID, error) {
	var id MemberID
	err := db.QueryRow("SELECT id FROM member WHERE PrimaryEmail = ?", u).Scan(&id)
	if err != nil && err == sql.ErrNoRows {
		err = fmt.Errorf("unknown primary email address")
		slog.Error(err.Error(), "username", u)
		return 0, err
	}
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}
	return id, nil
}
