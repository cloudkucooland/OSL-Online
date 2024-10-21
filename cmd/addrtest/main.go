package main

import (
	"context"
	"fmt"
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
	ids, err := model.ActiveMembers()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	for _, i := range ids {
		m, err := model.GetMember(i, true)
		if err != nil {
			panic(err)
		}

		if m.Address == "" { // a few of these exist
			continue
		}

		v, err := m.FormatAddress()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n\n", v)
	}
}
