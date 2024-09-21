package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/cloudkucooland/OSL-Online/model"
)

func main() {
	ctx, shutdown := context.WithCancel(context.Background())

	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB enviornment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}


	if err := model.Connect(ctx, dbpath); err != nil {
		slog.Error("startup", "message", "Error connecting to database", "error", err.Error())
		panic(err)
	}

	m, err := model.GetMember(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+x", m)

	shutdown()

}
