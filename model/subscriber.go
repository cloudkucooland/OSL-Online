package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type Subscriber struct {
	ID                int
	Name              string
	Attn              string
	Address           string
	AddressLine2      string
	City              string
	State             string
	Country           string
	PostalCode        string
	PrimaryPhone      string
	SecondaryPhone    string
	PrimaryEmail      string
	SecondaryEmail    string
	DateRecordCreated time.Time
	DatePaid          time.Time
	Doxology          string
	Newsletter        string
	Communication     string
}

type SubscriberImport struct {
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
	var n SubscriberImport

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
	return (&n).toSubscriber(), nil
}

func (n *SubscriberImport) toSubscriber() *Subscriber {
	return &Subscriber{
		ID:                n.ID,
		Name:              n.Name.String,
		Attn:              n.Attn.String,
		Address:           n.Address.String,
		AddressLine2:      n.AddressLine2.String,
		City:              n.City.String,
		State:             n.State.String,
		Country:           n.Country.String,
		PostalCode:        n.PostalCode.String,
		PrimaryPhone:      n.PrimaryPhone.String,
		SecondaryPhone:    n.SecondaryPhone.String,
		PrimaryEmail:      n.PrimaryEmail.String,
		SecondaryEmail:    n.SecondaryEmail.String,
		DateRecordCreated: n.DateRecordCreated,
		DatePaid:          n.DatePaid,
		Doxology:          n.Doxology.String,
		Newsletter:        n.Newsletter.String,
		Communication:     n.Communication.String,
	}
}

func (n *Subscriber) Store() error {
	_, err := db.Exec("REPLACE INTO subscriber (ID, Name, Attn, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, DateRecordCreated, DatePaid, Doxology, Newsletter, Communication) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", n.ID, n.Name, n.Attn, n.Address, n.AddressLine2, n.City, n.State, n.Country, n.PostalCode, n.PrimaryPhone, n.SecondaryPhone, n.PrimaryEmail, n.SecondaryEmail, n.DateRecordCreated, n.DatePaid, n.Doxology, n.Newsletter, n.Communication)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (n *SubscriberImport) Store() error {
	_, err := db.Exec("REPLACE INTO subscriber (ID, Name, Attn, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, DateRecordCreated, DatePaid, Doxology, Newsletter, Communication) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", n.ID, n.Name, n.Attn, n.Address, n.AddressLine2, n.City, n.State, n.Country, n.PostalCode, n.PrimaryPhone, n.SecondaryPhone, n.PrimaryEmail, n.SecondaryEmail, n.DateRecordCreated, n.DatePaid, n.Doxology, n.Newsletter, n.Communication)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func SetSubscriberField(id int, field string, value string) error {
	slog.Info("updating", "id", id, "field", field, "value", value)
	return nil
}
