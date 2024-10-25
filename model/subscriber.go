package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type SubscriberID int

type Subscriber struct {
	ID                SubscriberID
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
	FormattedAddr     string
}

type subNulls struct {
	ID                SubscriberID
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

func (id SubscriberID) Get() (*Subscriber, error) {
	var n subNulls

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
	s := (&n).toSubscriber()
	s.FormattedAddr, err = FormatAddress(s)
	if err != nil {
		slog.Error(err.Error())
		return s, err
	}
	return s, nil
}

func (n *subNulls) toSubscriber() *Subscriber {
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

func (n *subNulls) Store() error {
	_, err := db.Exec("REPLACE INTO subscriber (ID, Name, Attn, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, DateRecordCreated, DatePaid, Doxology, Newsletter, Communication) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", n.ID, n.Name, n.Attn, n.Address, n.AddressLine2, n.City, n.State, n.Country, n.PostalCode, n.PrimaryPhone, n.SecondaryPhone, n.PrimaryEmail, n.SecondaryEmail, n.DateRecordCreated, n.DatePaid, n.Doxology, n.Newsletter, n.Communication)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (id SubscriberID) SetField(field string, value string) error {
	slog.Info("updating", "id", id, "field", field, "value", value)

	if field == "id" {
		err := fmt.Errorf("cannot change ID")
		slog.Error(err.Error())
		return err
	}
	if strings.ContainsAny(field, "`;%") {
		err := fmt.Errorf("sql injection attempt [%s]", field)
		slog.Error(err.Error())
		return err
	}
	q := fmt.Sprintf("UPDATE `subscriber` SET `%s` = ? WHERE `id` = ?", field)

	switch field {
	case "DatePaid":
		value = strings.TrimSpace(value)
		if value == "" {
			value = "0001-01-01"
		}
		t, err := time.Parse("2006-01-02", value)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		if _, err := db.Exec(q, t, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	default:
		value = strings.TrimSpace(value)
		var ns sql.NullString
		if value == "" {
			ns.Valid = false
			ns.String = ""
		} else {
			ns.Valid = true
			ns.String = value
		}
		if _, err := db.Exec(q, ns, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	// XXX get ID
	if _, err := db.Exec("INSERT INTO auditlog VALUES (?, ?, ?, ?, CURRENT_DATE())", 0, id, field, value); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func (s *Subscriber) ISOCountry() string {
	return s.Country
}
