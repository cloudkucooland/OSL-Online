package main

import (
	"context"
	// "fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/cloudkucooland/OSL-Online/rest"
)

func main() {
	ctx, shutdown := context.WithCancel(context.Background())

	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB enviornment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}

	if err := model.Connect(ctx, dbpath); err != nil {
		slog.Error("startup", "message", "Error connecting to database", "error", err.Error())
		panic(err)
	}

	go rest.Start(ctx)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	sig := <-sigch
	slog.Info("shutdown", "requested by signal", sig)

	shutdown()
	time.Sleep(2)
}
