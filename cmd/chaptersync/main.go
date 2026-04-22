package main

import (
	"flag"
	"strconv"

	"github.com/cloudkucooland/OSL-Online/model"
)

func main() {
	flag.Parse()
	c, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		panic(err.Error())
	}
	id := model.ChapterID(c)

	model.NeedGAC()
	ctx, disconnect := model.ConnectCLI()
	defer disconnect()

	if err := id.ChapterSync(ctx); err != nil {
		panic(err.Error())
	}
}
