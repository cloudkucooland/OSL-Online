package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// MemberImport is the format used by the import tool and in the main query since NULLs are possible
type MemberImport struct {
	ID                  int
	MemberStatus        sql.NullString
	FirstName           sql.NullString
	MiddleName          sql.NullString
	LastName            sql.NullString
	PreferredName       sql.NullString
	Title               sql.NullString
	LifevowName         sql.NullString
	Suffix              sql.NullString
	Address             sql.NullString
	AddressLine2        sql.NullString
	City                sql.NullString
	State               sql.NullString
	Country             sql.NullString
	PostalCode          sql.NullString
	PrimaryPhone        sql.NullString
	SecondaryPhone      sql.NullString
	PrimaryEmail        sql.NullString
	SecondaryEmail      sql.NullString
	BirthDate           time.Time
	DateRecordCreated   time.Time
	Chapter             sql.NullString
	DateFirstVows       time.Time
	DateReaffirmation   time.Time
	DateRemoved         time.Time
	DateFirstProfession time.Time
	DateDeceased        time.Time
	DateNovitiate       time.Time
	Status              sql.NullString
	HowJoined           sql.NullString
	HowRemoved          sql.NullString
	ListInDirectory     sql.NullBool
	ListAddress         sql.NullBool
	ListPrimaryPhone    sql.NullBool
	ListSecondaryPhone  sql.NullBool
	ListPrimaryEmail    sql.NullBool
	ListSecondaryEmail  sql.NullBool
	Doxology            sql.NullString
	Newsletter          sql.NullString
	Communication       sql.NullString
	Occupation          sql.NullString
	Employeer           sql.NullString
	Denomination        sql.NullString
}

// Member is the format sent to the UI
type Member struct {
	ID                  int
	MemberStatus        string
	FirstName           string
	MiddleName          string
	LastName            string
	PreferredName       string
	Title               string
	LifevowName         string
	Suffix              string
	Address             string
	AddressLine2        string
	City                string
	State               string
	Country             string
	PostalCode          string
	PrimaryPhone        string
	SecondaryPhone      string
	PrimaryEmail        string
	SecondaryEmail      string
	BirthDate           time.Time
	DateRecordCreated   time.Time
	Chapter             string
	DateFirstVows       time.Time
	DateReaffirmation   time.Time
	DateRemoved         time.Time
	DateFirstProfession time.Time
	DateDeceased        time.Time
	DateNovitiate       time.Time
	Status              string
	HowJoined           string
	HowRemoved          string
	ListInDirectory     bool
	ListAddress         bool
	ListPrimaryPhone    bool
	ListSecondaryPhone  bool
	ListPrimaryEmail    bool
	ListSecondaryEmail  bool
	Doxology            string
	Newsletter          string
	Communication       string
	Occupation          string
	Employeer           string
	Denomination        string
}

// GetMember returns a populated Member struct, NULLs converted to ""
func GetMember(id int) (*Member, error) {
	var n MemberImport

	var bd, rc, fv, ra, dr, fp, dd, dn sql.NullString

	err := db.QueryRow("SELECT ID, MemberStatus, FirstName, MiddleName, LastName, PreferredName, Title, LifevowName, Suffix, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, BirthDate, DateRecordCreated, Chapter, DateFirstVows, DateReaffirmation, DateRemoved, DateFirstProfession, DateDeceased, DateNovitiate, Status, HowJoined, HowRemoved, ListInDirectory, ListAddress, ListPrimaryPhone, ListSecondaryPhone, ListPrimaryEmail, ListSecondaryEmail, Doxology, Newsletter, Communication, Occupation, Employeer, Denomination FROM member WHERE ID = ?", id).Scan(&n.ID, &n.MemberStatus, &n.FirstName, &n.MiddleName, &n.LastName, &n.PreferredName, &n.Title, &n.LifevowName, &n.Suffix, &n.Address, &n.AddressLine2, &n.City, &n.State, &n.Country, &n.PostalCode, &n.PrimaryPhone, &n.SecondaryPhone, &n.PrimaryEmail, &n.SecondaryEmail, &bd, &rc, &n.Chapter, &fv, &ra, &dr, &fp, &dd, &dn, &n.Status, &n.HowJoined, &n.HowRemoved, &n.ListInDirectory, &n.ListAddress, &n.ListPrimaryPhone, &n.ListSecondaryPhone, &n.ListPrimaryEmail, &n.ListSecondaryEmail, &n.Doxology, &n.Newsletter, &n.Communication, &n.Occupation, &n.Employeer, &n.Denomination)
	if err != nil && err == sql.ErrNoRows {
		err = fmt.Errorf("member not found")
		slog.Error(err.Error(), "id", id)
		return nil, err
	}
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	if bd.Valid {
		n.BirthDate, _ = time.Parse("2006-01-02", bd.String)
	}
	if rc.Valid {
		n.DateRecordCreated, _ = time.Parse("2006-01-02", rc.String)
	}
	if fv.Valid {
		n.DateFirstVows, _ = time.Parse("2006-01-02", fv.String)
	}
	if ra.Valid {
		n.DateReaffirmation, _ = time.Parse("2006-01-02", ra.String)
	}
	if dr.Valid {
		n.DateRemoved, _ = time.Parse("2006-01-02", dr.String)
	}
	if fp.Valid {
		n.DateFirstProfession, _ = time.Parse("2006-01-02", fp.String)
	}
	if dd.Valid {
		n.DateDeceased, _ = time.Parse("2006-01-02", dd.String)
	}
	if dn.Valid {
		n.DateNovitiate, _ = time.Parse("2006-01-02", dn.String)
	}

	return (&n).toMember(), nil
}

func (n *MemberImport) toMember() *Member {
	return &Member{
		ID:                  n.ID,
		MemberStatus:        n.MemberStatus.String,
		FirstName:           n.FirstName.String,
		MiddleName:          n.MiddleName.String,
		LastName:            n.LastName.String,
		PreferredName:       n.PreferredName.String,
		Title:               n.Title.String,
		LifevowName:         n.LifevowName.String,
		Suffix:              n.Suffix.String,
		Address:             n.Address.String,
		AddressLine2:        n.AddressLine2.String,
		City:                n.City.String,
		State:               n.State.String,
		Country:             n.Country.String,
		PostalCode:          n.PostalCode.String,
		PrimaryPhone:        n.PrimaryPhone.String,
		SecondaryPhone:      n.SecondaryPhone.String,
		PrimaryEmail:        n.PrimaryEmail.String,
		SecondaryEmail:      n.SecondaryEmail.String,
		BirthDate:           n.BirthDate,
		DateRecordCreated:   n.DateRecordCreated,
		Chapter:             n.Chapter.String,
		DateFirstVows:       n.DateFirstVows,
		DateReaffirmation:   n.DateReaffirmation,
		DateRemoved:         n.DateRemoved,
		DateFirstProfession: n.DateFirstProfession,
		DateDeceased:        n.DateDeceased,
		DateNovitiate:       n.DateNovitiate,
		Status:              n.Status.String,
		HowJoined:           n.HowJoined.String,
		HowRemoved:          n.HowRemoved.String,
		ListInDirectory:     n.ListInDirectory.Bool,
		ListAddress:         n.ListAddress.Bool,
		ListPrimaryPhone:    n.ListPrimaryPhone.Bool,
		ListSecondaryPhone:  n.ListSecondaryPhone.Bool,
		ListPrimaryEmail:    n.ListPrimaryEmail.Bool,
		ListSecondaryEmail:  n.ListSecondaryEmail.Bool,
		Doxology:            n.Doxology.String,
		Newsletter:          n.Newsletter.String,
		Communication:       n.Communication.String,
		Occupation:          n.Occupation.String,
		Employeer:           n.Employeer.String,
		Denomination:        n.Denomination.String,
	}
}

// probably don't want to use this one since it will write "" over the top of NULLs -- need a Member.toImport()
func (n *Member) Store() error {
	_, err := db.Exec("REPLACE INTO member (ID, MemberStatus, FirstName, MiddleName, LastName, PreferredName, Title, LifevowName, Suffix, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, BirthDate, DateRecordCreated, Chapter, DateFirstVows, DateReaffirmation, DateRemoved, DateFirstProfession, DateDeceased, DateNovitiate, Status, HowJoined, HowRemoved, ListInDirectory, ListAddress, ListPrimaryPhone, ListSecondaryPhone, ListPrimaryEmail, ListSecondaryEmail, Doxology, Newsletter, Communication, Occupation, Employeer, Denomination) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", n.ID, n.MemberStatus, n.FirstName, n.MiddleName, n.LastName, n.PreferredName, n.Title, n.LifevowName, n.Suffix, n.Address, n.AddressLine2, n.City, n.State, n.Country, n.PostalCode, n.PrimaryPhone, n.SecondaryPhone, n.PrimaryEmail, n.SecondaryEmail, n.BirthDate, n.DateRecordCreated, n.Chapter, n.DateFirstVows, n.DateReaffirmation, n.DateRemoved, n.DateFirstProfession, n.DateDeceased, n.DateNovitiate, n.Status, n.HowJoined, n.HowRemoved, n.ListInDirectory, n.ListAddress, n.ListPrimaryPhone, n.ListSecondaryPhone, n.ListPrimaryEmail, n.ListSecondaryEmail, n.Doxology, n.Newsletter, n.Communication, n.Occupation, n.Employeer, n.Denomination)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (n *MemberImport) Store() error {
	_, err := db.Exec("REPLACE INTO member (ID, MemberStatus, FirstName, MiddleName, LastName, PreferredName, Title, LifevowName, Suffix, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, BirthDate, DateRecordCreated, Chapter, DateFirstVows, DateReaffirmation, DateRemoved, DateFirstProfession, DateDeceased, DateNovitiate, Status, HowJoined, HowRemoved, ListInDirectory, ListAddress, ListPrimaryPhone, ListSecondaryPhone, ListPrimaryEmail, ListSecondaryEmail, Doxology, Newsletter, Communication, Occupation, Employeer, Denomination) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", n.ID, n.MemberStatus, n.FirstName, n.MiddleName, n.LastName, n.PreferredName, n.Title, n.LifevowName, n.Suffix, n.Address, n.AddressLine2, n.City, n.State, n.Country, n.PostalCode, n.PrimaryPhone, n.SecondaryPhone, n.PrimaryEmail, n.SecondaryEmail, n.BirthDate, n.DateRecordCreated, n.Chapter, n.DateFirstVows, n.DateReaffirmation, n.DateRemoved, n.DateFirstProfession, n.DateDeceased, n.DateNovitiate, n.Status, n.HowJoined, n.HowRemoved, n.ListInDirectory, n.ListAddress, n.ListPrimaryPhone, n.ListSecondaryPhone, n.ListPrimaryEmail, n.ListSecondaryEmail, n.Doxology, n.Newsletter, n.Communication, n.Occupation, n.Employeer, n.Denomination)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func SetMemberField(id int, field string, value string) error {
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
	q := fmt.Sprintf("UPDATE `member` SET `%s` = ? WHERE `id` = ?", field)
	slog.Info(q)

	switch field {
	case "ListInDirectory", "ListAddress", "ListPrimaryPhone", "ListSecondaryPhone", "ListPrimaryEmail", "ListSecondaryEmail":
		slog.Info("bool")
		var nb sql.NullBool
		nb.Valid = true
		nb.Bool = value == "true"
		if _, err := db.Exec(q, nb, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	case "BirthDate", "DateRecordCreated", "DateFirstVows", "DateReaffirmation", "DateRemoved", "DateFirstProfession", "DateDeceased", "DateNovitiate":
		slog.Info("date")
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
		slog.Info("string")
		var ns sql.NullString
		if value == "" || strings.TrimSpace(value) == "" {
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
