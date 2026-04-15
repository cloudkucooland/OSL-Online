package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/cloudkucooland/OSL-Online/model"
)

const usage = "auth username password level"

func main() {
	flag.Parse()
	username := model.Authname(flag.Arg(0))
	password := flag.Arg(1)
	if username == "" || password == "" {
		panic(usage)
	}

	ctx := context.WithValue(context.Background(), model.CtxKeyLevel, model.AuthLevelInternal)

	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB enviornment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}

	if err := model.Connect(ctx, dbpath); err != nil {
		slog.Error("startup", "message", "Error connecting to database", "error", err.Error())
		panic(err)
	}

	if err := username.SetAuthData(password); err != nil {
		panic(err.Error())
	}
}
