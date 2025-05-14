package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/cloudkucooland/OSL-Online/email"
	"github.com/cloudkucooland/OSL-Online/model"
)

func background(ctx context.Context) {
	ticker := time.NewTicker(time.Hour * 24)
	for {
		select {
		case <-ticker.C:
			slog.Info("daily tasks")
			doDaily()
		case <-ctx.Done():
			slog.Info("stopping background tasks")
			return
		}
	}
}

func doDaily() {
	now := time.Now()
	day := now.Day()
	month := now.Month()

	members, err := model.SearchBirthday(month, day)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	bdays := make([]*email.BirthdayEmailEntry, 0)
	for _, m := range members {
		bdays = append(bdays, &email.BirthdayEmailEntry{ID: int(m.ID), Name: m.OSLName()})
	}

	if err := email.SendBirthdayMail(bdays, month, day); err != nil {
		panic(err)
	}
}
