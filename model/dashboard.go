package model

import (
	// "database/sql"
	"fmt"
	"log/slog"
	"time"
)

// Dashboard is the format sent to the UI
type Dashboard_t struct {
	LifevowCount    int
	AnnualCount     int
	FriendCount     int
	SubscriberCount int
	ThisYearGiving  string
	LastYearGiving  string
}

func Dashboard() (Dashboard_t, error) {
	var d Dashboard_t

	rows, err := db.Query("SELECT MemberStatus, count(*) FROM member GROUP BY MemberStatus")
	if err != nil {
		slog.Error(err.Error())
		return d, err
	}
	defer rows.Close()
	for rows.Next() {
		var status string;
		var count int;
		if err = rows.Scan(&status, &count); err != nil {
			slog.Error(err.Error())
			continue
		}
		switch status {
		case "Annual Vows":
			d.AnnualCount = count
		case "Life Vows":
			d.LifevowCount = count
		case "Friend":
			d.FriendCount = count
		}
	}

	y := time.Now().Year()
	lastyear := fmt.Sprintf("%d-07-01", y - 1)
	thisyear := fmt.Sprintf("%d-07-01", y)
	nextyear := fmt.Sprintf("%d-07-01", y + 1)

	err = db.QueryRow("select sum(amount) from giving where date > ? and date < ?", lastyear, thisyear).Scan(&d.LastYearGiving)
	if err != nil {
		slog.Error(err.Error())
		return d, err
	}

	err = db.QueryRow("select sum(amount) from giving where date > ? and date < ?", thisyear, nextyear).Scan(&d.ThisYearGiving)
	if err != nil {
		slog.Error(err.Error())
		return d, err
	}

	return d, nil
}
