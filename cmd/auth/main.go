package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"strconv"

	"github.com/cloudkucooland/OSL-Online/model"
)

const usage = "auth username password level"

func main() {
	flag.Parse()
	username := flag.Arg(0)
	password := flag.Arg(1)
	level := flag.Arg(2)
	if username == "" || password == "" || level == "" {
		panic(usage)
	}

	ctx := context.Background()

	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB enviornment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}

	if err := model.Connect(ctx, dbpath); err != nil {
		slog.Error("startup", "message", "Error connecting to database", "error", err.Error())
		panic(err)
	}

	ll, err := strconv.Atoi(level)
	if err != nil {
		panic(err.Error())
	}

	if err := model.SetAuthData(username, password, ll); err != nil {
		panic(err.Error())
	}
}
