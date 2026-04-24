package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/cloudkucooland/OSL-Online/model"
)

var usauth string
var sgauth string

func main() {
	ctx, disconnect := model.ConnectCLI()
	defer disconnect()

	// zipfix2 needs a changer ID for the changelog
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

	ids, err := model.ActiveMemberIDs(ctx) // everybody
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
		time.Sleep(1 * time.Second)
	}
}
