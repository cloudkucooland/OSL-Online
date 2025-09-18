package main

import (
	"context"
	// "fmt"
	"log/slog"
	"os"
	// "strings"
	"time"

	"github.com/cloudkucooland/OSL-Online/model"
)

func main() {
	ctx := context.Background()

	bearer, err := getauth(ctx)
	if err != nil {
		panic(err)
	}

	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB enviornment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}

	if err := model.Connect(ctx, dbpath); err != nil {
		slog.Error("startup", "message", "Error connecting to database", "error", err.Error())
		panic(err)
	}

	ids, err := model.JustMemberIDsUS()
	// ids, err := model.FriendIDs()
	// ids, err := model.NecrologyIDs()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	for _, id := range ids {
		member, err := id.Get()
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		if member.Country != "US" {
			continue
		}
		err = getaddress(ctx, member, bearer)
		if err != nil {
			slog.Error(err.Error())
		}
		// default rate limit is 60/hr... yuk this is going to take DAYS
		time.Sleep(60 * time.Second)
	}
}
