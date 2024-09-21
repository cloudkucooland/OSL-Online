package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type Member struct {
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
	ListInDirectory     bool
	ListAddress         bool
	ListPrimaryPhone    bool
	ListSecondaryPhone  bool
	ListPrimaryEmail    bool
	ListSecondaryEmail  bool
	Doxology            sql.NullString
	Newsletter          sql.NullString
	Communication       sql.NullString
	Occupation          sql.NullString
	Employeer           sql.NullString
	Denomination        sql.NullString
}

func GetMember(id int) (*Member, error) {
	var m Member
	var n = &m

	// these will always be set, need NullString?
	var ynDirectory, ynAddress, ynPrimaryPhone, ynSecondaryPhone, ynPrimaryEmail, ynSecondaryEmail sql.NullString
	var bd, rc, fv, ra, dr, fp, dd, dn sql.NullString

	err := db.QueryRow("SELECT ID, MemberStatus, FirstName, MiddleName, LastName, PreferredName, Title, LifevowName, Suffix, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, BirthDate, DateRecordCreated, Chapter, DateFirstVows, DateReaffirmation, DateRemoved, DateFirstProfession, DateDeceased, DateNovitiate, Status, HowJoined, HowRemoved, ListInDirectory, ListAddress, ListPrimaryPhone, ListSecondaryPhone, ListPrimaryEmail, ListSecondaryEmail, Doxology, Newsletter, Communication, Occupation, Employeer, Denomination FROM member WHERE ID = ?", id).Scan(&n.ID, &n.MemberStatus, &n.FirstName, &n.MiddleName, &n.LastName, &n.PreferredName, &n.Title, &n.LifevowName, &n.Suffix, &n.Address, &n.AddressLine2, &n.City, &n.State, &n.Country, &n.PostalCode, &n.PrimaryPhone, &n.SecondaryPhone, &n.PrimaryEmail, &n.SecondaryEmail, &bd, &rc, &n.Chapter, &fv, &ra, &dr, &fp, &dd, &dn, &n.Status, &n.HowJoined, &n.HowRemoved, &ynDirectory, &ynAddress, &ynPrimaryPhone, &ynSecondaryPhone, &ynPrimaryEmail, &ynSecondaryEmail, &n.Doxology, &n.Newsletter, &n.Communication, &n.Occupation, &n.Employeer, &n.Denomination)
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

	n.ListInDirectory = yn2b(ynDirectory)
	n.ListAddress = yn2b(ynAddress)
	n.ListPrimaryPhone = yn2b(ynPrimaryPhone)
	n.ListSecondaryPhone = yn2b(ynSecondaryPhone)
	n.ListPrimaryEmail = yn2b(ynPrimaryEmail)
	n.ListSecondaryEmail = yn2b(ynSecondaryEmail)

	return n, nil
}

func SetMember(n *Member) error {
	_, err := db.Exec("REPLACE INTO member (ID, MemberStatus, FirstName, MiddleName, LastName, PreferredName, Title, LifevowName, Suffix, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, BirthDate, DateRecordCreated, Chapter, DateFirstVows, DateReaffirmation, DateRemoved, DateFirstProfession, DateDeceased, DateNovitiate, Status, HowJoined, HowRemoved, ListInDirectory, ListAddress, ListPrimaryPhone, ListSecondaryPhone, ListPrimaryEmail, ListSecondaryEmail, Doxology, Newsletter, Communication, Occupation, Employeer, Denomination) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", n.ID, n.MemberStatus, n.FirstName, n.MiddleName, n.LastName, n.PreferredName, n.Title, n.LifevowName, n.Suffix, n.Address, n.AddressLine2, n.City, n.State, n.Country, n.PostalCode, n.PrimaryPhone, n.SecondaryPhone, n.PrimaryEmail, n.SecondaryEmail, n.BirthDate, n.DateRecordCreated, n.Chapter, n.DateFirstVows, n.DateReaffirmation, n.DateRemoved, n.DateFirstProfession, n.DateDeceased, n.DateNovitiate, n.Status, n.HowJoined, n.HowRemoved, b2yn(n.ListInDirectory), b2yn(n.ListAddress), b2yn(n.ListPrimaryPhone), b2yn(n.ListSecondaryPhone), b2yn(n.ListPrimaryEmail), b2yn(n.ListSecondaryEmail), n.Doxology, n.Newsletter, n.Communication, n.Occupation, n.Employeer, n.Denomination)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func PrintMember(id int) {
	m, err := GetMember(id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Name: %s %s %s %s\n", m.Title.String, m.PreferredName.String, m.LastName.String, m.Suffix.String)
	fmt.Printf("Member Status: %s\n", m.MemberStatus.String)
	fmt.Println(" -- ")
}

func yn2b(yn sql.NullString) bool {
	if yn.Valid && yn.String == "YES" {
		return true
	}
	return false
}

func b2yn(b bool) sql.NullString {
	var s sql.NullString
	s.Valid = true
	if b {
		s.String = "YES"
	}
	s.String = "NO"
	return s
}
