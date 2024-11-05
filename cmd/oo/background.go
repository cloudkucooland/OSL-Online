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
			doDaily()
		case <-ctx.Done():
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
		bdays = append(bdays, &email.BirthdayEmailEntry{int(m.ID), m.OSLName()})
	}

	if err := email.SendBirthdayMail(bdays, month, day); err != nil {
		panic(err)
	}
}
