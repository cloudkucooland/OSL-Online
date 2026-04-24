package main

import (
	"flag"

	"github.com/cloudkucooland/OSL-Online/model"
)

const usage = "auth username password level"

func main() {
	flag.Parse()
	username := model.Authname(flag.Arg(0))
	password := flag.Arg(1)
	if username == "" || password == "" {
		panic(usage)
	}

	_, disconnect := model.ConnectCLI()
	defer disconnect()

	if err := username.SetAuthData(password); err != nil {
		panic(err.Error())
	}
}
