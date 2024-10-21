package model

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	// need a comment here to make lint happy
	_ "github.com/go-sql-driver/mysql"
)

// db is the private global used by all relevant functions to interact with the database
var db *sql.DB

// Connect tries to establish a connection to a MySQL/MariaDB database under the given URI and initializes the tables if they don"t exist yet.
func Connect(ctx context.Context, uri string) error {
	result, err := sql.Open("mysql", uri)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	db = result

	var version string
	if err := db.QueryRow("SELECT VERSION()").Scan(&version); err != nil {
		slog.Error(err.Error())
		return err
	}
	slog.Info("startup", "database", "connected", "version", version, "message", "connected to database")
	return nil
}

// Disconnect closes the database connection
// called only at server shutdown
func Disconnect() {
	slog.Info("shutdown", "message", "cleanly disconnected from database")
	if err := db.Close(); err != nil {
		slog.Error(err.Error())
	}
}

// makeNullString is used for values that may & might be inserted/updated as NULL in the database
func makeNullString(in interface{}) sql.NullString {
	var s string

	tmp, ok := in.(string)
	if ok {
		s = tmp
	} else {
		tmp, ok := in.(fmt.Stringer)
		if !ok {
			return sql.NullString{}
		}
		s = tmp.String()
	}
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func makeNullBool(in bool) sql.NullBool {
	return sql.NullBool{
		Bool:  in,
		Valid: true,
	}
}
