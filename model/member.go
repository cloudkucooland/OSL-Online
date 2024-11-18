package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type MemberID int

// memberNulls is the format used by the import tool and in the main query since NULLs are possible
type memberNulls struct {
	ID                 MemberID
	MemberStatus       sql.NullString
	FirstName          sql.NullString
	MiddleName         sql.NullString
	LastName           sql.NullString
	PreferredName      sql.NullString
	Title              sql.NullString
	LifevowName        sql.NullString
	Suffix             sql.NullString
	Address            sql.NullString
	AddressLine2       sql.NullString
	City               sql.NullString
	State              sql.NullString
	Country            sql.NullString
	PostalCode         sql.NullString
	PrimaryPhone       sql.NullString
	SecondaryPhone     sql.NullString
	PrimaryEmail       sql.NullString
	SecondaryEmail     sql.NullString
	BirthDate          sql.NullTime
	DateRecordCreated  sql.NullTime
	DateFirstVows      sql.NullTime
	DateReaffirmation  sql.NullTime
	DateRemoved        sql.NullTime
	DateDeceased       sql.NullTime
	DateNovitiate      sql.NullTime
	DateLifeVows       sql.NullTime
	Status             sql.NullString
	Leadership         sql.NullString
	HowJoined          sql.NullString
	HowRemoved         sql.NullString
	ListInDirectory    sql.NullBool
	ListAddress        sql.NullBool
	ListPrimaryPhone   sql.NullBool
	ListSecondaryPhone sql.NullBool
	ListPrimaryEmail   sql.NullBool
	ListSecondaryEmail sql.NullBool
	Doxology           sql.NullString
	Newsletter         sql.NullString
	Communication      sql.NullString
	Occupation         sql.NullString
	Employer           sql.NullString
	Denomination       sql.NullString
}

// Member is the format sent to the UI
type Member struct {
	ID                 MemberID
	MemberStatus       string
	FirstName          string
	MiddleName         string
	LastName           string
	PreferredName      string
	Title              string
	LifevowName        string
	Suffix             string
	Address            string
	AddressLine2       string
	City               string
	State              string
	Country            string
	PostalCode         string
	PrimaryPhone       string
	SecondaryPhone     string
	PrimaryEmail       string
	SecondaryEmail     string
	BirthDate          time.Time
	DateRecordCreated  time.Time
	DateFirstVows      time.Time
	DateReaffirmation  time.Time
	DateRemoved        time.Time
	DateDeceased       time.Time
	DateNovitiate      time.Time
	DateLifeVows       time.Time
	Status             string
	Leadership         string
	HowJoined          string
	HowRemoved         string
	ListInDirectory    bool
	ListAddress        bool
	ListPrimaryPhone   bool
	ListSecondaryPhone bool
	ListPrimaryEmail   bool
	ListSecondaryEmail bool
	Doxology           string
	Newsletter         string
	Communication      string
	Occupation         string
	Employer           string
	Denomination       string
	FormattedAddr      string // only populated for export to webui
}

// GetMember returns a populated Member struct, NULLs converted to ""
func (id MemberID) Get() (*Member, error) {
	n := &memberNulls{}

	err := db.QueryRow("SELECT ID, MemberStatus, FirstName, MiddleName, LastName, PreferredName, Title, LifevowName, Suffix, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, BirthDate, DateRecordCreated, DateFirstVows, DateReaffirmation, DateRemoved, DateDeceased, DateNovitiate, DateLifeVows, Status, Leadership, HowJoined, HowRemoved, ListInDirectory, ListAddress, ListPrimaryPhone, ListSecondaryPhone, ListPrimaryEmail, ListSecondaryEmail, Doxology, Newsletter, Communication, Occupation, Employer, Denomination FROM member WHERE ID = ?", id).Scan(&n.ID, &n.MemberStatus, &n.FirstName, &n.MiddleName, &n.LastName, &n.PreferredName, &n.Title, &n.LifevowName, &n.Suffix, &n.Address, &n.AddressLine2, &n.City, &n.State, &n.Country, &n.PostalCode, &n.PrimaryPhone, &n.SecondaryPhone, &n.PrimaryEmail, &n.SecondaryEmail, &n.BirthDate, &n.DateRecordCreated, &n.DateFirstVows, &n.DateReaffirmation, &n.DateRemoved, &n.DateDeceased, &n.DateNovitiate, &n.DateLifeVows, &n.Status, &n.Leadership, &n.HowJoined, &n.HowRemoved, &n.ListInDirectory, &n.ListAddress, &n.ListPrimaryPhone, &n.ListSecondaryPhone, &n.ListPrimaryEmail, &n.ListSecondaryEmail, &n.Doxology, &n.Newsletter, &n.Communication, &n.Occupation, &n.Employer, &n.Denomination)
	if err != nil && err == sql.ErrNoRows {
		err = fmt.Errorf("member not found")
		slog.Error(err.Error(), "id", id)
		return nil, err
	}
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	m := n.toMember()
	return m, nil
}

func (n *memberNulls) toMember() *Member {
	return &Member{
		ID:                 n.ID,
		MemberStatus:       n.MemberStatus.String,
		FirstName:          n.FirstName.String,
		MiddleName:         n.MiddleName.String,
		LastName:           n.LastName.String,
		PreferredName:      n.PreferredName.String,
		Title:              n.Title.String,
		LifevowName:        n.LifevowName.String,
		Suffix:             n.Suffix.String,
		Address:            n.Address.String,
		AddressLine2:       n.AddressLine2.String,
		City:               n.City.String,
		State:              n.State.String,
		Country:            n.Country.String,
		PostalCode:         n.PostalCode.String,
		PrimaryPhone:       n.PrimaryPhone.String,
		SecondaryPhone:     n.SecondaryPhone.String,
		PrimaryEmail:       n.PrimaryEmail.String,
		SecondaryEmail:     n.SecondaryEmail.String,
		BirthDate:          n.BirthDate.Time,
		DateRecordCreated:  n.DateRecordCreated.Time,
		DateFirstVows:      n.DateFirstVows.Time,
		DateReaffirmation:  n.DateReaffirmation.Time,
		DateRemoved:        n.DateRemoved.Time,
		DateDeceased:       n.DateDeceased.Time,
		DateNovitiate:      n.DateNovitiate.Time,
		DateLifeVows:       n.DateLifeVows.Time,
		Status:             n.Status.String,
		Leadership:         n.Leadership.String,
		HowJoined:          n.HowJoined.String,
		HowRemoved:         n.HowRemoved.String,
		ListInDirectory:    n.ListInDirectory.Bool,
		ListAddress:        n.ListAddress.Bool,
		ListPrimaryPhone:   n.ListPrimaryPhone.Bool,
		ListSecondaryPhone: n.ListSecondaryPhone.Bool,
		ListPrimaryEmail:   n.ListPrimaryEmail.Bool,
		ListSecondaryEmail: n.ListSecondaryEmail.Bool,
		Doxology:           n.Doxology.String,
		Newsletter:         n.Newsletter.String,
		Communication:      n.Communication.String,
		Occupation:         n.Occupation.String,
		Employer:           n.Employer.String,
		Denomination:       n.Denomination.String,
	}
}

func (n *Member) tomemberNulls() *memberNulls {
	return &memberNulls{
		ID:                 n.ID,
		MemberStatus:       makeNullString(n.MemberStatus),
		FirstName:          makeNullString(n.FirstName),
		MiddleName:         makeNullString(n.MiddleName),
		LastName:           makeNullString(n.LastName),
		PreferredName:      makeNullString(n.PreferredName),
		Title:              makeNullString(n.Title),
		LifevowName:        makeNullString(n.LifevowName),
		Suffix:             makeNullString(n.Suffix),
		Address:            makeNullString(n.Address),
		AddressLine2:       makeNullString(n.AddressLine2),
		City:               makeNullString(n.City),
		State:              makeNullString(n.State),
		Country:            makeNullString(n.Country),
		PostalCode:         makeNullString(n.PostalCode),
		PrimaryPhone:       makeNullString(n.PrimaryPhone),
		SecondaryPhone:     makeNullString(n.SecondaryPhone),
		PrimaryEmail:       makeNullString(n.PrimaryEmail),
		SecondaryEmail:     makeNullString(n.SecondaryEmail),
		BirthDate:          makeNullTime(n.BirthDate),
		DateRecordCreated:  makeNullTime(n.DateRecordCreated),
		DateFirstVows:      makeNullTime(n.DateFirstVows),
		DateReaffirmation:  makeNullTime(n.DateReaffirmation),
		DateRemoved:        makeNullTime(n.DateRemoved),
		DateDeceased:       makeNullTime(n.DateDeceased),
		DateNovitiate:      makeNullTime(n.DateNovitiate),
		DateLifeVows:       makeNullTime(n.DateLifeVows),
		Status:             makeNullString(n.Status),
		Leadership:         makeNullString(n.Leadership),
		HowJoined:          makeNullString(n.HowJoined),
		HowRemoved:         makeNullString(n.HowRemoved),
		ListInDirectory:    makeNullBool(n.ListInDirectory),
		ListAddress:        makeNullBool(n.ListAddress),
		ListPrimaryPhone:   makeNullBool(n.ListPrimaryPhone),
		ListSecondaryPhone: makeNullBool(n.ListSecondaryPhone),
		ListPrimaryEmail:   makeNullBool(n.ListPrimaryEmail),
		ListSecondaryEmail: makeNullBool(n.ListSecondaryEmail),
		Doxology:           makeNullString(n.Doxology),
		Newsletter:         makeNullString(n.Newsletter),
		Communication:      makeNullString(n.Communication),
		Occupation:         makeNullString(n.Occupation),
		Employer:           makeNullString(n.Employer),
		Denomination:       makeNullString(n.Denomination),
	}
}

func (n *memberNulls) Store() error {
	_, err := db.Exec("REPLACE INTO member (ID, MemberStatus, FirstName, MiddleName, LastName, PreferredName, Title, LifevowName, Suffix, Address, AddressLine2, City, State, Country, PostalCode, PrimaryPhone, SecondaryPhone, PrimaryEmail, SecondaryEmail, BirthDate, DateRecordCreated, DateFirstVows, DateReaffirmation, DateRemoved, DateDeceased, DateNovitiate, DateLifeVows, Status, Leadership, HowJoined, HowRemoved, ListInDirectory, ListAddress, ListPrimaryPhone, ListSecondaryPhone, ListPrimaryEmail, ListSecondaryEmail, Doxology, Newsletter, Communication, Occupation, Employer, Denomination) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", n.ID, n.MemberStatus, n.FirstName, n.MiddleName, n.LastName, n.PreferredName, n.Title, n.LifevowName, n.Suffix, n.Address, n.AddressLine2, n.City, n.State, n.Country, n.PostalCode, n.PrimaryPhone, n.SecondaryPhone, n.PrimaryEmail, n.SecondaryEmail, n.BirthDate, n.DateRecordCreated, n.DateFirstVows, n.DateReaffirmation, n.DateRemoved, n.DateDeceased, n.DateNovitiate, n.DateLifeVows, n.Status, n.Leadership, n.HowJoined, n.HowRemoved, n.ListInDirectory, n.ListAddress, n.ListPrimaryPhone, n.ListSecondaryPhone, n.ListPrimaryEmail, n.ListSecondaryEmail, n.Doxology, n.Newsletter, n.Communication, n.Occupation, n.Employer, n.Denomination)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (n *Member) Store() error {
	nn := n.tomemberNulls()
	return nn.Store()
}

func (id MemberID) SetMemberField(field string, value string, changer MemberID) error {
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

	switch field {
	case "ListInDirectory", "ListAddress", "ListPrimaryPhone", "ListSecondaryPhone", "ListPrimaryEmail", "ListSecondaryEmail":
		var nb sql.NullBool
		nb.Valid = true
		nb.Bool = value == "true"
		if _, err := db.Exec(q, nb, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	case "BirthDate", "DateRecordCreated", "DateFirstVows", "DateReaffirmation", "DateRemoved", "DateDeceased", "DateNovitiate", "DateLifeVows":
		value = strings.TrimSpace(value)
		if value == "" {
			value = zerotime
		}
		t, err := time.Parse(timeformat, value)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		if _, err := db.Exec(q, t, id); err != nil {
			slog.Error(err.Error())
			return err
		}
	case "MemberStatus":
		switch value {
		case "Removed":
			if _, err := db.Exec("UPDATE `member` SET `MemberStatus` = 'Removed', `ListInDirectory` = 0, `ListAddress` = 0, `ListPrimaryPhone` = 0, `ListSecondaryPhone` = 0, `ListPrimaryEmail` = 0, `ListSecondaryEmail` = 0, `Doxology` = 'none', `Newsletter` = 'none', `Communication` = 'none' WHERE id = ?", id); err != nil {
				slog.Error(err.Error())
				return err
			}
		case "Deceased":
			if _, err := db.Exec("UPDATE `member` SET `MemberStatus` = 'Deceased', `ListInDirectory` = 0, `ListAddress` = 0, `ListPrimaryPhone` = 0, `ListSecondaryPhone` = 0, `ListPrimaryEmail` = 0, `ListSecondaryEmail` = 0, `Doxology` = 'none', `Newsletter` = 'none', `Communication` = 'none' WHERE id = ?", id); err != nil {
				slog.Error(err.Error())
				return err
			}
		case "Annual Vows", "Friend", "Life Vows":
			if _, err := db.Exec(q, value, id); err != nil {
				slog.Error(err.Error())
				return err
			}
		default:
			err := fmt.Errorf("unknown MemberStatus")
			slog.Error(err.Error())
			return err
		}
	case "PrimaryPhone", "SecondaryPhone":
		var ns sql.NullString

		m, err := id.Get()
		if err != nil {
			return err
		}

		pn, err := FormatPhoneNumber(value, m.Country)
		if err != nil {
			return err
		}

		if pn == "" {
			ns.Valid = false
			ns.String = ""
		} else {
			ns.Valid = true
			ns.String = pn
		}

		if _, err := db.Exec(q, ns, id); err != nil {
			slog.Error(err.Error())
			return err
		}
		value = pn // for logging
	default:
		var ns sql.NullString
		value = strings.TrimSpace(value)
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

	id.ChangeLogStore(ChangeLogEntry{
		Changer: changer,
		Field:   field,
		Value:   value,
	})

	return nil
}

func Create(firstname, lastname string) (MemberID, error) {
	if firstname == "" || lastname == "" {
		err := fmt.Errorf("name cannot be null")
		slog.Error(err.Error())
		return 0, err
	}

	n := memberNulls{
		ID:                0,
		MemberStatus:      sql.NullString{Valid: true, String: "Friend"},
		FirstName:         sql.NullString{Valid: true, String: firstname},
		LastName:          sql.NullString{Valid: true, String: lastname},
		DateRecordCreated: sql.NullTime{Valid: true, Time: time.Now()},
	}

	res, err := db.Exec("INSERT INTO member (MemberStatus, FirstName, LastName, DateRecordCreated) VALUES (?,?,?,?)", n.MemberStatus, n.FirstName, n.LastName, n.DateRecordCreated)
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}
	last, err := res.LastInsertId()
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}
	return MemberID(last), nil
}

func (n *Member) CleanUnlisted() {
	zt, _ := time.Parse(timeformat, zerotime)

	if !n.ListInDirectory  {
		n.FirstName= ""
		n.LastName= ""
		n.MiddleName= ""
		n.PreferredName= ""
		n.Title  = ""
		n.LifevowName  = ""
		n.Suffix  = ""
		n.BirthDate  = zt
		n.DateNovitiate  = zt
		n.DateRemoved  = zt
		n.DateFirstVows  = zt
		n.DateReaffirmation  = zt
		n.DateDeceased  = zt
		n.Status  = ""
		n.Occupation  = ""
		n.Employer  = ""
		n.Denomination  = ""
		n.HowRemoved  = ""
		n.ListAddress  = false
		n.ListPrimaryPhone  = false
		n.ListSecondaryPhone  = false
		n.ListPrimaryEmail  = false
		n.ListSecondaryEmail  = false
	}

	if !n.ListAddress  {
		n.Address  = ""
		n.AddressLine2  = ""
		n.City  = ""
		n.State  = ""
		n.Country  = ""
		n.PostalCode  = ""
	}

	if !n.ListPrimaryPhone  {
		n.PrimaryPhone  = ""
	}

	if !n.ListSecondaryPhone  {
		n.SecondaryPhone  = ""
	}

	if !n.ListPrimaryEmail  {
		n.PrimaryEmail  = ""
	}

	if !n.ListSecondaryEmail  {
		n.SecondaryEmail  = ""
	}
}

func (n *Member) OSLName() string {
	name := n.OSLShortName()

	name += " " + n.LastName
	if n.Suffix != "" {
		name += " " + n.Suffix
	}

	if n.MemberStatus == "Annual Vows" || n.MemberStatus == "Life Vows" {
		name += ", OSL"
	}
	return name
}

func (n *Member) OSLShortName() string {
	var name string
	firstname := false
	if n.MemberStatus == "Annual Vows" || n.MemberStatus == "Life Vows" {
		name = n.Title + " "
	}
	if n.MemberStatus == "Life Vows" && n.LifevowName != "" {
		name += n.LifevowName
		firstname = true
	}
	if !firstname && n.PreferredName != "" {
		name += n.PreferredName
		firstname = true
	}
	if !firstname && n.FirstName != "" {
		name += n.FirstName
	}
	return name
}

func (m *Member) SetChapters(chapters ...int) error {
	if _, err := db.Exec("DELETE FROM `chaptermembers` WHERE `member` = ?", m.ID); err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, ch := range chapters {
		_, err := db.Exec("INSERT INTO `chaptermembers` (`chapter`, `member`) VALUES (?, ?)", ch, m.ID)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	return nil
}

func (m *Member) GetChapters() ([]int, error) {
	rows, err := db.Query("SELECT `chapter` FROM `chaptermembers` WHERE `member` = ?", m.ID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	chapters := make([]int, 0)
	var ch int
	for rows.Next() {
		if err = rows.Scan(&ch); err != nil {
			slog.Error(err.Error())
			continue
		}
		chapters = append(chapters, ch)
	}
	return chapters, nil
}

func (m Member) ISOCountry() string {
	return m.Country
}
