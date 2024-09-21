package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type Subscriber struct {
	ID                int
	Name              sql.NullString
	Attn              sql.NullString
	Address           sql.NullString
	AddressLine2      sql.NullString
	City              sql.NullString
	State             sql.NullString
	Country           sql.NullString
	PostalCode        sql.NullString
	PrimaryPhone      sql.NullString
	SecondaryPhone    sql.NullString
	PrimaryEmail      sql.NullString
	SecondaryEmail    sql.NullString
	DateRecordCreated time.Time
	DatePaid          time.Time
	Doxology          sql.NullString
	Newsletter        sql.NullString
	Communication     sql.NullString
}

func GetSubscriber(id int) (*Subscriber, error) {
	var n Subscriber

	var created, paid sql.NullString

	err := db.QueryRow("SELECT ID, Name, Attn, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, DateRecordCreated, DatePaid, Doxology, Newsletter, Communication FROM subscriber WHERE ID = ?", id).Scan(&n.ID, &n.Name, &n.Attn, &n.Address, &n.AddressLine2, &n.City, &n.State, &n.Country, &n.PostalCode, &n.PrimaryPhone, &n.SecondaryPhone, &n.PrimaryEmail, &n.SecondaryEmail, &created, &paid, &n.Doxology, &n.Newsletter, &n.Communication)
	if err != nil && err == sql.ErrNoRows {
		err = fmt.Errorf("subscriber not found")
		slog.Error(err.Error(), "id", id)
		return nil, err
	}
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	if created.Valid {
		n.DateRecordCreated, _ = time.Parse("2006-01-02", created.String)
	}
	if paid.Valid {
		n.DatePaid, _ = time.Parse("2006-01-02", paid.String)
	}
	return &n, nil
}

func SetSubscriber(n *Subscriber) error {
	_, err := db.Exec("REPLACE INTO subscriber (ID, Name, Attn, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, DateRecordCreated, DatePaid, Doxology, Newsletter, Communication) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", n.ID, n.Name, n.Attn, n.Address, n.AddressLine2, n.City, n.State, n.Country, n.PostalCode, n.PrimaryPhone, n.SecondaryPhone, n.PrimaryEmail, n.SecondaryEmail, n.DateRecordCreated, n.DatePaid, n.Doxology, n.Newsletter, n.Communication)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
