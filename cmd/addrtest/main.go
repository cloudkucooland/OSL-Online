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

	// get the IDs of all active members, friends and benefactors
	ids, err := model.ActiveMemberIDs()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	for _, i := range ids {
		if i == 0 {
			continue
		}
		m, err := i.Get(true)
		if err != nil {
			panic(err)
		}
		formatted, err := model.FormatAddress(m)
		if err != nil {
			slog.Info(err.Error())
		}
		fmt.Printf("%+v\n\n", formatted)

		if m.Address != "" && formatted == "" {
			panic(m.OSLName())
		}
	}

	sids, err := model.ActiveSubscriberIDs()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	for _, i := range sids {
		m, err := i.Get()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n\n", m.FormattedAddr)

		if m.FormattedAddr == "" {
			panic(m.Name)
		}
	}
}
