package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudkucooland/OSL-Online/model"
)

var errTime time.Time

func main() {
	ctx, shutdown := context.WithCancel(context.Background())

	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB not set")
	}

	if err := model.Connect(ctx, dbpath); err != nil {
		panic(err)
	}

	f, err := os.Open("export.csv")
	if err != nil {
		panic(err)
	}

	r := csv.NewReader(f)
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
	r.ReuseRecord = true

	errTime, _ = time.Parse("1/2/2006", "1/1/0001")
	id := 1000 // starting ID
	for {
		d, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		// the header will err
		checkID, err := strconv.Atoi(d[102])
		if err != nil || checkID == 0 {
			continue
		}
		// we don't use the ID because it's crazy huge
		id = id + 1

		if d[1] == "ORGANIZATION" || d[14] == "ORGANIZATION" || d[109] == "Subscriber" || d[109] == "Doxology" || d[109] == "Doxology Only" || d[109] == "Copy Editor" || d[109] == "Author" || d[109] == "Contributor" {
			doOrg(id, d)
		} else {
			doMember(id, d)
		}
	}
	shutdown()
}

func doMember(id int, d []string) {
	var m model.MemberImport

	m.ID = model.MemberID(id)
	m.LastName.Valid = d[0] != ""
	m.LastName.String = d[0]
	m.FirstName.Valid = d[1] != ""
	m.FirstName.String = d[1]
	m.PrimaryPhone.Valid = d[2] != ""
	m.PrimaryPhone.String = d[2]
	m.Address.Valid = d[4] != ""
	m.Address.String = d[4]
	m.ListAddress.Valid = true
	m.ListAddress.Bool = d[5] == "No" // Unlisted vs. listed
	t, err := time.Parse("1/2/2006", d[20])
	if err != nil {
		m.BirthDate = errTime
	} else {
		m.BirthDate = t
	}
	m.SecondaryPhone.Valid = d[25] != ""
	m.SecondaryPhone.String = d[25]
	m.ListSecondaryPhone.Valid = true
	m.ListSecondaryPhone.Bool = d[26] == "No" // Unlisted vs. listed
	m.City.Valid = d[27] != ""
	m.City.String = d[27]
	m.Country.Valid = d[31] != ""
	m.Country.String = d[31]
	t, err = time.Parse("1/2/2006", d[34])
	if err != nil {
		m.DateRecordCreated = errTime
	} else {
		m.DateRecordCreated = t
	}
	// very few, only new people have this set, use it as fallback for firstvow
	t, err = time.Parse("1/2/2006", d[35])
	if err != nil {
		m.DateFirstVows = errTime
	} else {
		m.DateFirstVows = t
	}

	t, err = time.Parse("1/2/2006", d[38])
	if err != nil {
		m.DateRemoved = errTime
	} else {
		m.DateRemoved = t
	}

	t, err = time.Parse("1/2/2006", d[39])
	if err != nil {
		m.DateDeceased = errTime
	} else {
		m.DateDeceased = t
	}

	m.PrimaryEmail.Valid = d[42] != ""
	m.PrimaryEmail.String = d[42]
	m.ListPrimaryEmail.Valid = true
	m.ListPrimaryEmail.Bool = d[44] == "No" // Unlisted vs. listed

	m.SecondaryEmail.Valid = d[45] != ""
	m.SecondaryEmail.String = d[45]
	m.ListSecondaryEmail.Valid = true
	m.ListSecondaryEmail.Bool = d[46] == "No" // Unlisted vs. listed

	m.Employer.Valid = d[51] != ""
	m.Employer.String = d[51]

	t, err = time.Parse("1/2/2006", d[92])
	if err != nil {
		m.DateFirstVows = errTime
	} else {
		m.DateFirstVows = t
	}

	m.HowJoined.Valid = d[97] != ""
	m.HowJoined.String = d[97]

	m.HowRemoved.Valid = d[98] != ""
	m.HowRemoved.String = d[98]

	m.ListInDirectory.Valid = true
	m.ListInDirectory.Bool = d[101] == "Yes"

	m.MemberStatus.Valid = true
	m.Doxology.Valid = true
	m.Newsletter.Valid = true
	m.Communication.Valid = true
	switch d[109] {
	case "Annual Vows", "":
		m.Doxology.String = "electronic"
		m.Newsletter.String = "electronic"
		m.Communication.String = "mailed"
		m.MemberStatus.String = "Annual Vows"
	case "Life Vows":
		m.Doxology.String = "electronic"
		m.Newsletter.String = "electronic"
		m.Communication.String = "mailed"
		m.MemberStatus.String = d[109]
	case "Hold Periodicals - Life Vows":
		m.Doxology.String = "electronic"
		m.Newsletter.String = "electronic"
		m.Communication.String = "mailed"
		m.MemberStatus.String = "Life Vows"
	case "Hold Periodicals - Annual Vows":
		m.Doxology.String = "none"
		m.Newsletter.String = "electronic"
		m.Communication.String = "mailed"
		m.MemberStatus.String = "Annual Vows"
	case "Removed":
		m.Doxology.String = "none"
		m.Newsletter.String = "none"
		m.Communication.String = "none"
		m.MemberStatus.String = "Removed"
	case "_Contributor":
		m.Doxology.String = "none"
		m.Newsletter.String = "none"
		m.Communication.String = "mailed"
		m.MemberStatus.String = "Contributor"
	default:
		fmt.Printf("Unknown member status: [%s]\n", d[109])
		m.Doxology.String = "none"
		m.Newsletter.String = "none"
		m.Communication.String = "none"
		m.MemberStatus.String = "Removed"
	}
	m.MemberStatus.Valid = true

	m.MiddleName.Valid = d[110] != ""
	m.MiddleName.String = d[110]

	switch d[111] {
	case "Emailed", "":
		m.Newsletter.String = "electronic"
	case "No Font":
		m.Newsletter.String = "none"
	case "Mailed":
		m.Newsletter.String = "mailed"
	default:
		fmt.Printf("Unknown Newsletter setting: [%s]\n", d[111])
	}

	t, err = time.Parse("1/2/2006", d[112])
	if err != nil {
		m.DateNovitiate = errTime
	} else {
		m.DateNovitiate = t
	}

	m.Occupation.Valid = d[113] != ""
	m.Occupation.String = d[113]

	m.ListPrimaryPhone.Valid = true
	m.ListPrimaryPhone.Bool = d[116] == "No" // Unlisted vs. listed

	switch d[119] {
	case "Clergy", "Deacon", "Chaplain":
		m.Status.String = "clergy"
		m.Status.Valid = true
	case "Lay", "Laity", "", "Military":
		m.Status.String = "laity"
		m.Status.Valid = true
	case "Student":
		m.Status.String = "student"
		m.Status.Valid = true
	case "Retired":
		m.Status.String = "retired laity"
		m.Status.Valid = true
	case "Retired Clergy":
		m.Status.String = "retired clergy"
		m.Status.Valid = true
	default:
		fmt.Printf("\nUnknown status: [%s]\n", d[119])
		m.Status.Valid = false
	}

	m.PreferredName.Valid = d[120] != ""
	m.PreferredName.String = d[120]

	t, err = time.Parse("1/2/2006", d[128])
	if err != nil {
		m.DateReaffirmation = errTime
	} else {
		m.DateReaffirmation = t
	}

	m.PreferredName.Valid = d[129] != ""
	m.PreferredName.String = d[129]

	m.State.Valid = d[137] != ""
	m.State.String = d[137]

	if d[138] != "" && d[138] != "OSL" {
		m.Suffix.Valid = true
		m.Suffix.String = strings.TrimSuffix(d[138], ", OSL")
	} else {
		m.Suffix.Valid = false
		m.Suffix.String = ""
	}

	switch d[141] {
	case "Br.", "Br":
		m.Title.String = "Br."
		m.Title.Valid = true
	case "Sr.", "Sr":
		m.Title.String = "Sr."
		m.Title.Valid = true
	case "Sibling", "Sib":
		m.Title.String = "Sibling"
		m.Title.Valid = true
	case "Rev.", "Mr.", "Ms.", "Mrs.", "Organization", "Padre", "":
		m.Title.Valid = false
		m.Title.String = ""
	default:
		fmt.Printf("\nUnknown title [%s]\n", d[141])
		m.Title.Valid = false
		m.Title.String = ""
	}

	if m.MemberStatus.String == "Removed" {
		m.Title.Valid = false
		m.Title.String = ""
	}

	m.PostalCode.Valid = d[184] != ""
	m.PostalCode.String = d[184]

	// fmt.Printf("%+v\n%+v\n", d, m)
	(&m).Store()
	// model.PrintMember(m.ID)
}

func doOrg(id int, d []string) {
	var s model.SubscriberImport

	s.ID = model.SubscriberID(id)
	s.Name.Valid = d[0] != ""
	s.Name.String = d[0]
	s.Attn.Valid = d[1] != "" && d[1] != "ORGANIZATION"
	s.Attn.String = d[1]
	s.PrimaryPhone.Valid = d[2] != ""
	s.PrimaryPhone.String = d[2]
	s.Address.Valid = d[4] != ""
	s.Address.String = d[4]
	s.SecondaryPhone.Valid = d[25] != ""
	s.SecondaryPhone.String = d[25]
	s.City.Valid = d[27] != ""
	s.City.String = d[27]
	s.Country.Valid = d[31] != ""
	s.Country.String = d[31]
	t, err := time.Parse("1/2/2006", d[34])
	if err != nil {
		s.DateRecordCreated = errTime
	} else {
		s.DateRecordCreated = t
	}

	s.PrimaryEmail.Valid = d[42] != ""
	s.PrimaryEmail.String = d[42]

	s.SecondaryEmail.Valid = d[45] != ""
	s.SecondaryEmail.String = d[45]

	s.Doxology.Valid = true
	s.Doxology.String = "mailed"
	s.Newsletter.Valid = true
	s.Newsletter.String = "mailed"
	s.Communication.Valid = true
	s.Communication.String = "none"

	t, err = time.Parse("1/2/2006", d[128])
	if err != nil {
		s.DatePaid = errTime
	} else {
		s.DatePaid = t
	}

	s.State.Valid = d[137] != ""
	s.State.String = d[137]

	s.PostalCode.Valid = d[184] != ""
	s.PostalCode.String = d[184]

	(&s).Store()
	// processed, _ := model.GetSubscriber(id)
	// fmt.Println(processed)
}
