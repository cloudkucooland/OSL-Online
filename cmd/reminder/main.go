package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/cloudkucooland/OSL-Online/email"
	"github.com/cloudkucooland/OSL-Online/model"
)

func main() {
	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB enviornment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}

	if err := model.Connect(context.Background(), dbpath); err != nil {
		slog.Error("startup", "message", "Error connecting to database", "error", err.Error())
		panic(err)
	}

	members, err := model.ReminderAnnual()
	if err != nil {
		panic(err)
	}

	for _, id := range members {
		if err := email.SendReminder(id); err != nil {
			panic(err)
		}
	}
}
