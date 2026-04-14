package main

import (
	"context"
	"flag"
	"os"
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

	ctx := context.WithValue(context.Background(), model.CtxKeyLevel, model.AuthLevelInternal)

	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB enviornment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}

	if err := model.Connect(ctx, dbpath); err != nil {
		panic(err)
	}

	gac := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if gac == "" {
		panic("GOOGLE_APPLICATION_CREDENTIALS enviornment var not set.")
	}

	if err := id.ChapterSync(ctx); err != nil {
		panic(err.Error())
	}
}
