package model

import (
	"context"
	"log/slog"
	"os"
)

// ConnectCLI is a helper for CLI tools to connect to the database and setup a context with internal auth level.
// It panics on error because it's intended for simple CLI tools.
func ConnectCLI() (context.Context, func()) {
	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB environment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}

	ctx := context.WithValue(context.Background(), CtxKeyLevel, AuthLevelInternal)
	if err := Connect(ctx, dbpath); err != nil {
		slog.Error("startup", "message", "Error connecting to database", "error", err.Error())
		panic(err)
	}

	return ctx, Disconnect
}

// NeedGAC ensures GOOGLE_APPLICATION_CREDENTIALS is set, or panics.
func NeedGAC() {
	gac := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if gac == "" {
		panic("GOOGLE_APPLICATION_CREDENTIALS environment var not set.")
	}
}
