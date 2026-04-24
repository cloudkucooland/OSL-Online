package main

import (
	"github.com/cloudkucooland/OSL-Online/model"
)

func main() {
	model.NeedGAC()
	ctx, disconnect := model.ConnectCLI()
	defer disconnect()

	if err := model.DoxologySync(ctx); err != nil {
		panic(err.Error())
	}
}
