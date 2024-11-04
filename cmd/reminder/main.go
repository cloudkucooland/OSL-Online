package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/cloudkucooland/OSL-Online/model"
	// "github.com/cloudkucooland/OSL-Online/email"
)

func main() {
	ctx := context.Background()

	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB enviornment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}

	if err := model.Connect(ctx, dbpath); err != nil {
		slog.Error("startup", "message", "Error connecting to database", "error", err.Error())
		panic(err)
	}

	members, err := model.ReminderAnnual()
	if err != nil {
		panic(err)
	}

	for _, id := range members {
		fmt.Println(id)
		/* if err := email.SendReminder(id); err != nil {
			panic(err)
		} */
	}
}
