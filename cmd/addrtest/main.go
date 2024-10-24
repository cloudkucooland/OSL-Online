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
		m, err := i.Get(true)
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

		/* if m.Address == "" { // a few of these exist
			continue
		} */

		v, err := m.FormatAddress()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n\n", v)
	}
}
