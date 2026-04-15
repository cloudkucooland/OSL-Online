package model

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/sethvargo/go-password/password"
)

type AuthLevel uint8

const (
	AuthLevelView     AuthLevel = iota // view public data (member/default)
	AuthLevelFullView                  // read-only superuser (prior/canon)
	AuthLevelManager                   // change/add members (council)
	AuthLevelAdmin                     // full system control (elected)
	AuthLevelInternal                  // internal system use
)

type Authname string

func (u Authname) String() string {
	return string(u)
}

func LevelFromContext(ctx context.Context) AuthLevel {
	level, ok := ctx.Value(CtxKeyLevel).(AuthLevel)
	if !ok {
		return AuthLevelView
	}
	return level
}

func IDFromContext(ctx context.Context) (MemberID, error) {
	id, ok := ctx.Value(CtxKeyID).(MemberID)
	if !ok {
		slog.Error("no id found in context")
		return 0, fmt.Errorf("no id found in context")
	}
	return id, nil
}

func (u Authname) getAuthData(ctx context.Context) (string, AuthLevel, error) {
	var pwhash string
	var leadership string

	query := `SELECT a.pwhash, m.Leadership FROM auth a JOIN member m ON a.user = m.PrimaryEmail WHERE a.user = ?`

	err := db.QueryRowContext(ctx, query, u).Scan(&pwhash, &leadership)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, fmt.Errorf("user %s not found", u)
		}
		slog.Error("database error in getAuthData", "err", err)
		return "", 0, err
	}

	var level AuthLevel
	switch leadership {
	case "elected":
		level = AuthLevelAdmin
	case "council":
		level = AuthLevelManager
	case "prior", "canon":
		level = AuthLevelFullView
	default:
		level = AuthLevelView
	}

	return pwhash, level, nil
}

func (u Authname) Authenticate(ctx context.Context, password string) (AuthLevel, error) {
	pwhash, level, err := u.getAuthData(ctx)
	if err != nil {
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pwhash), []byte(password)); err != nil {
		slog.Error("login failed", "username", u, "err", err)
		return 0, err
	}

	return level, nil
}

// SetAuthData no longer needs a level passed in; it only manages the credential
func (u Authname) SetAuthData(pw string) error {
	slog.Info("updating password", "username", u)
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	if err != nil {
		return err
	}

	_, err = db.Exec("REPLACE INTO auth (user, pwhash) VALUES (?,?)", u, bytes)
	if err != nil {
		return err
	}
	return nil
}

func (u Authname) Register(ctx context.Context) (string, error) {
	slog.Info("registering user", "username", u)

	// Ensure they exist in the member table first
	if _, err := u.GetID(ctx); err != nil {
		return "", err
	}

	newPassword, err := password.Generate(10, 3, 0, false, true)
	if err != nil {
		slog.Error("password generation failed", "err", err)
		return "", err
	}

	if err := u.SetAuthData(newPassword); err != nil {
		slog.Error("failed to set auth data", "err", err)
		return "", err
	}

	return newPassword, nil
}

func (u Authname) GetID(ctx context.Context) (MemberID, error) {
	var id MemberID
	err := db.QueryRowContext(ctx, "SELECT id FROM member WHERE PrimaryEmail = ?", u).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("unknown primary email address")
		}
		slog.Error("database error in GetID", "err", err)
		return 0, err
	}
	return id, nil
}
