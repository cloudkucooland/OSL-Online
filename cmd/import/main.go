package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	// "log/slog"
	"os"
	"strconv"
	"time"

	"github.com/cloudkucooland/OSL-Online/model"
)

func main() {
	ctx, shutdown := context.WithCancel(context.Background())
	var m model.Member

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

	errTime, _ := time.Parse("1/2/2006", "1/1/1800")

	for {
		d, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		// the header will err
		m.ID, err = strconv.Atoi(d[102])
		if err != nil {
			// fmt.Println(err.Error())
			continue
		}

		/* skip organizations for now
		if d[1] == "ORGANIZATION" {
			continue
		} */

		m.LastName.Valid = d[0] != ""
		m.LastName.String = d[0]
		m.FirstName.Valid = d[1] != ""
		m.FirstName.String = d[1]
		m.PrimaryPhone.Valid = d[2] != ""
		m.PrimaryPhone.String = d[2]
		m.Address.Valid = d[4] != ""
		m.Address.String = d[4]
		m.ListAddress = d[5] != "No" // Unlisted vs. listed
		t, err := time.Parse("1/2/2006", d[20])
		if err != nil {
			m.BirthDate = errTime
		} else {
			m.BirthDate = t
		}
		m.SecondaryPhone.Valid = d[25] != ""
		m.SecondaryPhone.String = d[25]
		m.ListSecondaryPhone = d[26] != "No" // Unlisted vs. listed
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
		m.ListPrimaryEmail = d[44] != "No" // Unlisted vs. listed

		m.SecondaryEmail.Valid = d[45] != ""
		m.SecondaryEmail.String = d[45]
		m.ListSecondaryEmail = d[46] != "No" // Unlisted vs. listed

		m.Employeer.Valid = d[51] != ""
		m.Employeer.String = d[51]

		t, err = time.Parse("1/2/2006", d[92])
		if err != nil {
			m.DateFirstVows = errTime
		} else {
			m.DateFirstVows = t
		}

		t, err = time.Parse("1/2/2006", d[93])
		if err != nil {
			m.DateFirstProfession = errTime
		} else {
			m.DateFirstProfession = t
		}

		m.HowJoined.Valid = d[97] != ""
		m.HowJoined.String = d[97]

		m.HowRemoved.Valid = d[98] != ""
		m.HowRemoved.String = d[98]

		m.ListInDirectory = d[101] == "Yes"

		// `MemberStatus` enum('Annual Vows','Life Vows','Subscriber','Contributor','Removed') NOT NULL,
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
		case "Subscriber":
			m.Doxology.String = "mailed"
			m.Newsletter.String = "none"
			m.Communication.String = "mailed"
			m.MemberStatus.String = d[109]
		case "Contributor", "Copy Editor", "Author":
			m.Doxology.String = "none"
			m.Newsletter.String = "none"
			m.Communication.String = "none"
			m.MemberStatus.String = "Contributor"
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
		case "Doxology Only", "Doxology":
			m.Doxology.String = "mailed"
			m.Newsletter.String = "none"
			m.Communication.String = "none"
			m.MemberStatus.String = "Subscriber"
		default:
			fmt.Printf("Unknown member status: [%s]\n", d[109])
			m.Doxology.String = "electronic"
			m.Newsletter.String = "electronic"
			m.Communication.String = "mailed"
			m.MemberStatus.String = "Annual Vows"
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

		m.ListPrimaryPhone = d[116] != "No" // Unlisted vs. listed

		switch d[119] {
		case "Clergy", "Deacon", "Chaplain":
			m.Status.String = "clergy"
		case "Lay", "Laity", "", "Military":
			m.Status.String = "laity"
		case "Student":
			m.Status.String = "student"
		case "Retired":
			m.Status.String = "retired laity"
		case "Retired Clergy":
			m.Status.String = "retired clergy"
		default:
			// fmt.Printf("\nUnknown status: [%s]\n", d[119])
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
			m.Suffix.String = d[138]
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

		model.SetMember(&m)

		model.PrintMember(m.ID)
	}
	shutdown()
}
