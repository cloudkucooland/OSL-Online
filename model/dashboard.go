package model

import (
	// "database/sql"
	"fmt"
	"log/slog"
	"time"
)

// Dashboard is the format sent to the UI
type Dashboard_t struct {
	LifevowCount         int
	AnnualCount          int
	FriendCount          int
	SubscriberCount      int
	ThisYearGiving       string
	LastYearGiving       string
	AnnualVowsWhoGave    int
	LifeVowsWhoGave      int
	AnnualVowsReaffirmed int
	LifeVowsCheckin      int
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
		var status string
		var count int
		if err = rows.Scan(&status, &count); err != nil {
			slog.Error(err.Error())
			continue
		}
		switch status {
		case ANNUAL:
			d.AnnualCount = count
		case LIFE:
			d.LifevowCount = count
		case FRIEND:
			d.FriendCount = count
		}
	}

	y := time.Now().Year()
	lastyear := fmt.Sprintf("%d-07-01", y-1)
	thisyear := fmt.Sprintf("%d-07-01", y)
	nextyear := fmt.Sprintf("%d-07-01", y+1)

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

	err = db.QueryRow("select count(*) from subscriber where DatePaid > DATE_SUB(CURRENT_DATE(), INTERVAL 366 DAY)").Scan(&d.SubscriberCount)
	if err != nil {
		slog.Error(err.Error())
		return d, err
	}

	countrows, err := db.Query("select memberstatus, count(*) from member where id in (select distinct id from giving where date > ?) group by memberstatus", thisyear)
	if err != nil {
		slog.Error(err.Error())
		return d, err
	}
	defer countrows.Close()
	for countrows.Next() {
		var status string
		var count int
		if err = countrows.Scan(&status, &count); err != nil {
			slog.Error(err.Error())
			continue
		}
		switch status {
		case ANNUAL:
			d.AnnualVowsWhoGave = count
		case LIFE:
			d.LifeVowsWhoGave = count
		case FRIEND:
			d.SubscriberCount++
		}
	}

	reaffirmed, err := db.Query("select memberstatus, count(*) from member where DateReaffirmation > ? group by memberstatus", thisyear)
	if err != nil {
		slog.Error(err.Error())
		return d, err
	}
	defer reaffirmed.Close()
	for reaffirmed.Next() {
		var status string
		var count int
		if err = reaffirmed.Scan(&status, &count); err != nil {
			slog.Error(err.Error())
			continue
		}
		switch status {
		case ANNUAL:
			d.AnnualVowsReaffirmed = count
		case LIFE:
			d.LifeVowsCheckin = count
		}
	}

	return d, nil
}
