package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/cloudkucooland/OSL-Online/model"
)

var usauth string
var sgauth string

func main() {
	ctx := context.WithValue(context.Background(), model.CtxKeyLevel, model.AuthLevelInternal)
	ctx = context.WithValue(ctx, model.CtxKeyID, model.MemberID(0))

	var err error

	usauth, err = getauth(ctx)
	if err != nil {
		panic(err)
	}

	sgauth, err = getOneMapToken()
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

	ids, err := model.ActiveMemberIDs(ctx) // everybody
	// ids, err := model.JustMemberIDsUS(ctx)
	// ids, err := model.FriendIDs(ctx)
	// ids, err := model.NecrologyIDs(ctx)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	for _, id := range ids {
		member, err := id.Get(ctx)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		err = getaddress(ctx, member)
		if err != nil {
			slog.Error(err.Error())
		}
		// default rate limit is 60/hr... yuk this is going to take DAYS
		// time.Sleep(60 * time.Second)
		time.Sleep(1 * time.Second)
	}
}
