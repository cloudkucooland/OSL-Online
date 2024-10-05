package model

import (
	"database/sql"
	"log/slog"
	"time"
)

func ReportNotRenewed() ([]*Member, error) {
	var members []*Member
	var n MemberImport

	var ra sql.NullString

	rows, err := db.Query("SELECT FirstName, LastName, PreferredName, Title, Address, AddressLine2, City, State, Country, PostalCode, PrimaryEmail, DateReaffirmation FROM member WHERE MemberStatus = 'Annual Vows' AND DateReaffirmation < DATE_SUB(CURRENT_DATE(), INTERVAL 365 DAY) ORDER BY DateReaffirmation")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&n.FirstName, &n.LastName, &n.PreferredName, &n.Title, &n.Address, &n.AddressLine2, &n.City, &n.State, &n.Country, &n.PostalCode, &n.PrimaryEmail, &ra)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		if ra.Valid {
			n.DateReaffirmation, _ = time.Parse("2006-01-02", ra.String)
		}
		members = append(members, (&n).toMember())
	}
	return members, nil
}

func ReportExpired() ([]*Member, error) {
	var members []*Member
	var n MemberImport

	var ra sql.NullString

	rows, err := db.Query("SELECT FirstName, LastName, PreferredName, Title, Address, AddressLine2, City, State, Country, PostalCode, PrimaryEmail, DateReaffirmation FROM member WHERE MemberStatus = 'Annual Vows' AND DateReaffirmation < DATE_SUB(CURRENT_DATE(), INTERVAL 730 DAY) ORDER BY DateReaffirmation")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&n.FirstName, &n.LastName, &n.PreferredName, &n.Title, &n.Address, &n.AddressLine2, &n.City, &n.State, &n.Country, &n.PostalCode, &n.PrimaryEmail, &ra)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		if ra.Valid {
			n.DateReaffirmation, _ = time.Parse("2006-01-02", ra.String)
		}
		members = append(members, (&n).toMember())
	}
	return members, nil
}

func ReportEmail() ([]*Member, error) {
	var members []*Member
	var n MemberImport

	rows, err := db.Query("SELECT MemberStatus, FirstName, LastName, PreferredName, Title, LifevowName, Suffix, PrimaryEmail, SecondaryEmail, ListPrimaryEmail, ListSecondaryEmail, Doxology, Newsletter, Communication FROM member WHERE MemberStatus = 'Annual Vows' OR MemberStatus = 'Life Vows' ORDER BY LastName, FirstName")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&n.MemberStatus, &n.FirstName, &n.LastName, &n.PreferredName, &n.Title, &n.LifevowName, &n.Suffix, &n.PrimaryEmail, &n.SecondaryEmail, &n.ListPrimaryEmail, &n.ListSecondaryEmail, &n.Doxology, &n.Newsletter, &n.Communication)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		members = append(members, (&n).toMember())
	}
	return members, nil
}

func ReportAnnual() ([]*Member, error) {
	var members []*Member
	var n MemberImport

	rows, err := db.Query("SELECT FirstName, LastName, PreferredName, Title, Suffix, Address, AddressLine2, City, State, Country, PostalCode, Doxology, Newsletter, Communication FROM member WHERE MemberStatus = 'Annual Vows' ORDER BY LastName, FirstName")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&n.FirstName, &n.LastName, &n.PreferredName, &n.Title, &n.Suffix, &n.Address, &n.AddressLine2, &n.City, &n.State, &n.Country, &n.PostalCode, &n.Doxology, &n.Newsletter, &n.Communication)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		members = append(members, (&n).toMember())
	}
	return members, nil
}

func ReportLife() ([]*Member, error) {
	var members []*Member
	var n MemberImport

	rows, err := db.Query("SELECT FirstName, LastName, PreferredName, Title, Suffix, LifevowName, Address, AddressLine2, City, State, Country, PostalCode, Doxology, Newsletter, Communication FROM member WHERE MemberStatus = 'Life Vows' ORDER BY LastName, FirstName")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&n.FirstName, &n.LastName, &n.PreferredName, &n.Title, &n.Suffix, &n.LifevowName, &n.Address, &n.AddressLine2, &n.City, &n.State, &n.Country, &n.PostalCode, &n.Doxology, &n.Newsletter, &n.Communication)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		members = append(members, (&n).toMember())
	}
	return members, nil
}

func ReportSubscriber() ([]*Subscriber, error) {
	var subscribers []*Subscriber
	var n SubscriberImport

	rows, err := db.Query("SELECT Name, Attn, Address, AddressLine2, City, State, Country, PostalCode, Doxology, Newsletter, Communication, DatePaid FROM subscriber ORDER BY Name")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&n.Name, &n.Attn, &n.Address, &n.AddressLine2, &n.City, &n.State, &n.Country, &n.PostalCode, &n.Doxology, &n.Newsletter, &n.Communication, &n.DatePaid)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		subscribers = append(subscribers, (&n).toSubscriber())
	}
	return subscribers, nil
}

// Returns a slice of IDs
func ActiveMembers() ([]int, error) {
	var id int
	list := make([]int, 0, 1000)

	rows, err := db.Query("SELECT id FROM member WHERE MemberStatus != 'Removed'")
	if err != nil {
		slog.Error(err.Error())
		return list, err
	}

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			slog.Error(err.Error())
			return list, err
		}
		list = append(list, id)
	}
	return list, nil
}
