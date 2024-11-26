package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/cloudkucooland/OSL-Online/model"
)

func main() {
	ctx := context.Background()

	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB enviornment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}

	if err := model.Connect(ctx, dbpath); err != nil {
		slog.Error("startup", "message", "Error connecting to database", "error", err.Error())
		panic(err)
	}

	if err := model.Friendzone(); err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
