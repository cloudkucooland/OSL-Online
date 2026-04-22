//go:build OMIT
// +build OMIT

package main

import (
	"encoding/csv"
	"os"

	"github.com/cloudkucooland/OSL-Online/model"
)

func main() {
	ctx, disconnect := model.ConnectCLI()
	defer disconnect()

	members, err := model.NewMemberIDs(ctx)
	if err != nil {
		panic(err)
	}

	r := csv.NewWriter(os.Stdout)
	_ = r.Write([]string{"DateFirstVows", "OSLName"})

	for _, id := range members {
		m, err := id.Get(ctx)
		if err != nil {
			continue
		}

		_ = r.Write([]string{m.DateFirstVows.Format("January 2, 2006"), m.OSLName()})
	}
	r.Flush()
}
