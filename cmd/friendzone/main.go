package main

import (
	"log/slog"

	"github.com/cloudkucooland/OSL-Online/model"
)

func main() {
	ctx, disconnect := model.ConnectCLI()
	defer disconnect()

	if err := model.Friendzone(ctx); err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
