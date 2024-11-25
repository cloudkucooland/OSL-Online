package model

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"time"
)

// reportMemberQuery returns full member data, unlisted data is included
func reportMemberQuery(query string) ([]*Member, error) {
	members := make([]*Member, 0)

	rows, err := db.Query(query)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id MemberID
		if err := rows.Scan(&id); err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		member, _ := id.Get()
		members = append(members, member)
	}
	return members, nil
}

func ReportExpired(w io.Writer) error {
	members, err := reportMemberQuery("SELECT id FROM member WHERE MemberStatus = 'Annual Vows' AND DateReaffirmation < DATE_SUB(CURRENT_DATE(), INTERVAL 730 DAY) ORDER BY DateReaffirmation")
	if err != nil {
		return err
	}

	r := csv.NewWriter(w)
	_ = r.Write([]string{"DateReaffirmation", "OSLName", "FormattedAddr", "PrimaryEmail"})

	for _, m := range members {
		f, err := FormatAddress(m)
		if err != nil {
			continue
		}
		_ = r.Write([]string{m.DateReaffirmation.Format(time.DateOnly), m.OSLName(), f, m.PrimaryEmail})
	}
	r.Flush()
	return nil
}

func ReportAnnual(w io.Writer) error {
	members, err := reportMemberQuery("SELECT id FROM member WHERE MemberStatus = 'Annual Vows' ORDER BY LastName, FirstName")
	if err != nil {
		return err
	}

	r := csv.NewWriter(w)
	_ = r.Write([]string{"OSLName", "OSLShortName", "DateReaffirmation", "FormattedAddress", "PrimaryEmail"})

	for _, m := range members {
		f, err := FormatAddress(m)
		if err != nil {
			continue
		}
		_ = r.Write([]string{m.OSLName(), m.OSLShortName(), m.DateReaffirmation.Format(time.DateOnly), f, m.PrimaryEmail})
	}
	r.Flush()
	return nil
}

func ReportLife(w io.Writer) error {
	members, err := reportMemberQuery("SELECT id FROM member WHERE MemberStatus = 'Life Vows' ORDER BY LastName, FirstName")
	if err != nil {
		return err
	}

	r := csv.NewWriter(w)
	_ = r.Write([]string{"OSLName", "OSLShortName", "FormattedAddress", "PrimaryEmail"})

	for _, m := range members {
		f, err := FormatAddress(m)
		if err != nil {
			continue
		}
		_ = r.Write([]string{m.OSLName(), m.OSLShortName(), f, m.PrimaryEmail})
	}
	r.Flush()
	return nil
}

func ReportAllEmail(w io.Writer) error {
	r := csv.NewWriter(w)
	_ = r.Write([]string{"OSLName", "OSLShortName", "MemberStatus", "PrimaryEmail", "Address"})

	m, err := ActiveMemberIDs()
	if err != nil {
		return err
	}

	for _, id := range m {
		n, err := id.Get()
		if err != nil {
			slog.Error(err.Error())
			err = nil
			continue
		}
		f, err := FormatAddress(n)
		if err != nil {
			continue
		}
		_ = r.Write([]string{n.OSLName(), n.OSLShortName(), n.MemberStatus, n.PrimaryEmail, f})
	}
	r.Flush()
	return nil
}

func ReportBarb(w io.Writer) error {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec", "X"}
	r := csv.NewWriter(w)
	_ = r.Write([]string{"OSLName", "MemberStatus", "PrimaryEmail", "Month"})

	m, err := ActiveMemberIDs()
	if err != nil {
		return err
	}
	permonth := int(len(m) / 12)

	month := 0
	i := 0
	for _, id := range m {
		if i > permonth {
			i = 0
			month = month + 1
		}
		n, err := id.Get()
		if err != nil {
			slog.Error(err.Error())
			err = nil
			continue
		}

		_ = r.Write([]string{n.OSLName(), n.MemberStatus, n.PrimaryEmail, months[month]})
		i = i + 1
	}
	r.Flush()
	return nil
}

// ReportFontEmail writes a report structured for Google Groups CSV upload
func ReportFontEmailed(w io.Writer) error {
	r := csv.NewWriter(w)
	_ = r.Write([]string{"Group Email [Required]", "Member Email", "Member Type", "Member Role"})

	rows, err := db.Query("SELECT PrimaryEmail FROM member WHERE MemberStatus != 'Removed' AND MemberStatus != 'Deceased' AND PrimaryEmail IS NOT NULL AND Newsletter = 'electronic' ORDER BY LastName, FirstName")
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var e string
		if err := rows.Scan(&e); err != nil {
			slog.Error(err.Error())
			// return err
			continue
		}
		_ = r.Write([]string{"font@saint-luke.net", e, "USER", "MEMBER"})
	}
	r.Flush()
	return nil
}

/* func ReportSubscriber() ([]*Subscriber, error) {
	subscribers := make([]*Subscriber, 0)

	rows, err := db.Query("SELECT id FROM subscriber ORDER BY Name")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id SubscriberID
		if err := rows.Scan(&id); err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		sub, _ := id.Get()
		subscribers = append(subscribers, sub)
	}
	return subscribers, nil
} */

// ActiveMemberIDs returns All Annual Vows, Life Vows and Friends
func ActiveMemberIDs() ([]MemberID, error) {
	return reportMemberIDQuery("SELECT id FROM member WHERE MemberStatus IN ('Annual Vows', 'Life Vows', 'Friend') ORDER BY LastName, FirstName")
}

// ActiveMemberIDsUS returns All Annual Vows, Life Vows and Friends in the US
func ActiveMemberIDsUS() ([]MemberID, error) {
	return reportMemberIDQuery("SELECT id FROM member WHERE MemberStatus IN ('Annual Vows', 'Life Vows', 'Friend') AND Country = 'US' ORDER BY LastName, FirstName")
}

// AnnualMemberIDs does what it says
func AnnualMemberIDs() ([]MemberID, error) {
	return reportMemberIDQuery("SELECT id FROM member WHERE MemberStatus = 'Annual Vows' ORDER BY LastName, FirstName")
}

// LifeMemberIDs does what it says
func LifeMemberIDs() ([]MemberID, error) {
	return reportMemberIDQuery("SELECT id FROM member WHERE MemberStatus = 'Life Vows' ORDER BY LastName, FirstName")
}

// NewMemberIDs does what it says
func NewMemberIDs() ([]MemberID, error) {
	return reportMemberIDQuery("SELECT id FROM member WHERE MemberStatus = 'Annual Vows' AND DateFirstVows > DATE_SUB(CURRENT_DATE(), INTERVAL 730 DAY) ORDER BY LastName, FirstName")
}

// FriendsIDs does what it says
func FriendIDs() ([]MemberID, error) {
	return reportMemberIDQuery("SELECT id FROM member WHERE MemberStatus = 'Friend' ORDER BY LastName, FirstName")
}

// ReminderAnnual returns those who have not reaffirmed in the past year
func ReminderAnnual() ([]MemberID, error) {
	return reportMemberIDQuery("SELECT id FROM member WHERE MemberStatus = 'Annual Vows' AND DateReaffirmation < DATE_SUB(CURRENT_DATE(), INTERVAL 365 DAY) ORDER BY LastName, FirstName")
}

func TestMemberIDs() ([]MemberID, error) {
	return []MemberID{1078}, nil
}

func reportMemberIDQuery(query string) ([]MemberID, error) {
	list := make([]MemberID, 0, 500)

	rows, err := db.Query(query)
	if err != nil {
		slog.Error(err.Error())
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			slog.Error(err.Error())
			return list, err
		}
		list = append(list, MemberID(id))
	}
	return list, nil
}

// Returns a slice of IDs
/* func ActiveSubscriberIDs() ([]SubscriberID, error) {
	return reportSubscriberIDQuery("SELECT id FROM subscriber WHERE DatePaid > DATE_SUB(CURRENT_DATE(), INTERVAL 366 DAY) ORDER BY Name")
} */

func ReportAllSubscribers(w io.Writer) error {
	r := csv.NewWriter(w)
	_ = r.Write([]string{"ID", "Name", "Attn", "Address", "AddressLine2", "City", "State", "Country", "PostalCode", "PrimaryPhone", "SecondaryPhone", "PrimaryEmail", "SecondaryEmail", "DateRecordCreated", "DatePaid", "Doxology", "Newsletter", "Communication"})

	subscribers, err := reportSubscriberIDQuery("SELECT id FROM subscriber ORDER BY id")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, id := range subscribers {
		s, err := id.Get()
		if err != nil {
			continue
		}
		subscriber := []string{fmt.Sprintf("%d", s.ID), s.Name, s.Attn, s.Address, s.AddressLine2, s.City, s.State, s.Country, s.PostalCode, s.PrimaryPhone, s.SecondaryPhone, s.PrimaryEmail, s.SecondaryEmail, s.DateRecordCreated.Format(time.DateOnly), s.DatePaid.Format(time.DateOnly), s.Doxology, s.Newsletter, s.Communication}
		_ = r.Write(subscriber)
	}
	r.Flush()
	return nil
}

func reportSubscriberIDQuery(query string) ([]SubscriberID, error) {
	list := make([]SubscriberID, 0, 50)

	rows, err := db.Query(query)
	if err != nil {
		slog.Error(err.Error())
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			slog.Error(err.Error())
			return list, err
		}
		list = append(list, SubscriberID(id))
	}
	return list, nil
}

func ReportAvery(w io.Writer) error {
	members := make([]addressFormatter, 0)

	ids, err := ActiveMemberIDsUS()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, id := range ids {
		m, err := id.Get()
		if err != nil {
			slog.Error(err.Error())
			continue
			// return err
		}

		members = append(members, m)
	}
	AveryLabels(w, members)
	return nil
}

func DoxologyPrinted(w io.Writer) error {
	r := csv.NewWriter(w)
	_ = r.Write([]string{"Last Name", "First Name", "Address", "City", "State", "Zip Code", "Country"})

	members, err := reportMemberIDQuery("SELECT id FROM member WHERE MemberStatus != 'Removed' AND MemberStatus != 'Deceased' AND DateReaffirmation > DATE_SUB(CURRENT_DATE(), INTERVAL 366 DAY) AND Doxology = 'mailed' ORDER BY LastName, FirstName")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, id := range members {
		m, err := id.Get()
		if err != nil {
			continue
		}
		member := []string{m.LastName, m.FirstName, m.Address, m.City, m.State, m.PostalCode, m.Country}
		_ = r.Write(member)
	}

	subscribers, err := reportSubscriberIDQuery("SELECT id FROM subscriber WHERE DatePaid > DATE_SUB(CURRENT_DATE(), INTERVAL 366 DAY) AND Doxology = 'mailed' ORDER BY Name")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, id := range subscribers {
		s, err := id.Get()
		if err != nil {
			continue
		}
		subscriber := []string{s.Name, s.Attn, s.Address, s.City, s.State, s.PostalCode, s.Country}
		_ = r.Write(subscriber)
	}
	r.Flush()
	return nil
}

// DoxologyEmailed writes a report that is structured for Google Groups CSV upload
func DoxologyEmailed(w io.Writer) error {
	r := csv.NewWriter(w)
	_ = r.Write([]string{"Group Email [Required]", "Member Email", "Member Type", "Member Role", "Name"})

	members, err := reportMemberIDQuery("SELECT id FROM member WHERE MemberStatus != 'Removed' AND Doxology = 'electronic' AND PrimaryEmail IS NOT NULL ORDER BY LastName, FirstName")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, id := range members {
		m, err := id.Get()
		if err != nil {
			continue
		}
		_ = r.Write([]string{"doxology@saint-luke.net", m.PrimaryEmail, "USER", "MEMBER", m.OSLName()})
	}

	subscribers, err := reportSubscriberIDQuery("SELECT id FROM subscriber WHERE Doxology = 'electronic' AND PrimaryEmail IS NOT NULL AND DatePaid > DATE_SUB(CURRENT_DATE(), INTERVAL 730 DAY) ORDER BY Name")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, id := range subscribers {
		s, err := id.Get()
		if err != nil {
			continue
		}
		_ = r.Write([]string{"doxology@saint-luke.net", s.PrimaryEmail, "USER", "MEMBER", s.Name + " : " + s.Attn})
	}
	r.Flush()
	return nil
}
