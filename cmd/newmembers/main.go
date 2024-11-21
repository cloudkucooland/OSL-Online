//go:build OMIT
// +build OMIT

package main

import (
	"context"
	"encoding/csv"
	"log/slog"
	"os"
	// "time"

	"github.com/cloudkucooland/OSL-Online/model"
)

func main() {
	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB enviornment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}

	if err := model.Connect(context.Background(), dbpath); err != nil {
		slog.Error("startup", "message", "Error connecting to database", "error", err.Error())
		panic(err)
	}

	members, err := model.NewMemberIDs()
	if err != nil {
		panic(err)
	}

	r := csv.NewWriter(os.Stdout)
	_ = r.Write([]string{"DateFirstVows", "OSLName"})

	for _, id := range members {
		m, err := id.Get()
		if err != nil {
			continue
		}

		_ = r.Write([]string{m.DateFirstVows.Format("January 2, 2006"), m.OSLName()})
	}
	r.Flush()
}
