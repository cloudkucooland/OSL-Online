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

	todaybdays, err := model.SearchBirthday(month, day)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if len(todaybdays) > 0 {
		bdayemails := make([]*email.BirthdayEmailEntry, 0)
		for _, m := range todaybdays {
			bdayemails = append(bdayemails, &email.BirthdayEmailEntry{ID: int(m.ID), Name: m.OSLName()})
		}

		if err := email.SendBirthdayMail(bdayemails, month, day); err != nil {
			panic(err)
		}
	}
}
