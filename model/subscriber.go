package model

import (
	"context"
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
	Doxology          communicationPref
	Newsletter        communicationPref
	Communication     communicationPref
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
	DateRecordCreated sql.NullTime
	DatePaid          sql.NullTime
	Doxology          sql.NullString
	Newsletter        sql.NullString
	Communication     sql.NullString
}

func (id SubscriberID) Get(ctx context.Context) (*Subscriber, error) {
	var n subNulls

	err := db.QueryRowContext(ctx, "SELECT ID, Name, Attn, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, DateRecordCreated, DatePaid, Doxology, Newsletter, Communication FROM subscriber WHERE ID = ?", id).Scan(&n.ID, &n.Name, &n.Attn, &n.Address, &n.AddressLine2, &n.City, &n.State, &n.Country, &n.PostalCode, &n.PrimaryPhone, &n.SecondaryPhone, &n.PrimaryEmail, &n.SecondaryEmail, &n.DateRecordCreated, &n.DatePaid, &n.Doxology, &n.Newsletter, &n.Communication)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Warn("subscriber not found", "id", id)
			return nil, fmt.Errorf("subscriber %d not found", id)
		}
		slog.Error("database error in Subscriber.Get", "err", err, "id", id)
		return nil, fmt.Errorf("database error: %w", err)
	}

	s := (&n).toSubscriber()
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
		DateRecordCreated: n.DateRecordCreated.Time,
		DatePaid:          n.DatePaid.Time,
		Doxology:          communicationPref(n.Doxology.String),
		Newsletter:        communicationPref(n.Newsletter.String),
		Communication:     communicationPref(n.Communication.String),
	}
}

/* func (n *Subscriber) Store() error {
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
} */

func (id SubscriberID) SetField(ctx context.Context, field string, value string) error {
	changer, _ := ctx.Value(CtxKeyID).(MemberID)
	slog.Info("updating subscriber field", "id", id, "field", field, "value", value, "by", changer)

	if field == "id" {
		slog.Warn("attempt to change subscriber ID", "id", id, "by", changer)
		return fmt.Errorf("cannot change ID")
	}
	if strings.ContainsAny(field, "`;%") {
		slog.Warn("potential sql injection attempt in SetField", "field", field, "by", changer)
		return fmt.Errorf("invalid field name [%s]", field)
	}
	q := fmt.Sprintf("UPDATE `subscriber` SET `%s` = ? WHERE `id` = ?", field)

	switch field {
	case "DatePaid":
		value = strings.TrimSpace(value)
		if value == "" {
			value = zerotime
		}
		t, err := time.Parse(timeformat, value)
		if err != nil {
			slog.Warn("invalid date format in SetField", "value", value, "id", id)
			return fmt.Errorf("invalid date format: %w", err)
		}
		if _, err := db.ExecContext(ctx, q, t, id); err != nil {
			slog.Error("database error in SetField (DatePaid)", "err", err, "id", id)
			return fmt.Errorf("database error: %w", err)
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
		if _, err := db.ExecContext(ctx, q, ns, id); err != nil {
			slog.Error("database error in SetField", "err", err, "id", id, "field", field)
			return fmt.Errorf("database error: %w", err)
		}
	}

	return nil
}

func (s *Subscriber) ISOCountry() string {
	return s.Country
}
