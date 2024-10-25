package model

import (
	"encoding/csv"
	"io"
	"log/slog"
	"time"
)

func reportMemberQuery(query string) ([]*Member, error) {
	var members []*Member

	rows, err := db.Query(query)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		var id MemberID
		if err := rows.Scan(&id); err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		member, _ := id.Get(true)
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
	r.Write([]string{"DateReaffirmation", "OSLName", "FormattedAddr", "PrimaryEmail"})

	for _, m := range members {
		r.Write([]string{m.DateReaffirmation.Format(time.DateOnly), m.OSLName(), m.FormattedAddr, m.PrimaryEmail})
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
	r.Write([]string{"OSLName", "OSLShortName", "DateReaffirmation", "FormattedAddress", "PrimaryEmail"})

	for _, m := range members {
		r.Write([]string{m.OSLName(), m.OSLShortName(), m.DateReaffirmation.Format(time.DateOnly), m.FormattedAddr, m.PrimaryEmail})
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
	r.Write([]string{"OSLName", "OSLShortName", "FormattedAddress", "PrimaryEmail"})

	for _, m := range members {
		r.Write([]string{m.OSLName(), m.OSLShortName(), m.FormattedAddr, m.PrimaryEmail})
	}
	r.Flush()
	return nil
}

func ReportAllEmail(w io.Writer) error {
	r := csv.NewWriter(w)
	r.Write([]string{"OSLName", "OSLShortName", "MemberStatus", "PrimaryEmail"})

	m, err := ActiveMemberIDs()
	if err != nil {
		return err
	}

	for _, id := range m {
		n, err := id.Get(true)
		if err != nil {
			slog.Error(err.Error())
			err = nil
			continue
		}
		r.Write([]string{n.OSLName(), n.OSLShortName(), n.MemberStatus, n.PrimaryEmail})
	}
	r.Flush()
	return nil
}

// structured for Google Groups CSV upload
func ReportFontEmailed(w io.Writer) error {
	r := csv.NewWriter(w)
	r.Write([]string{"Group Email [Required]", "Member Email", "Member Type", "Member Role"})

	rows, err := db.Query("SELECT PrimaryEmail FROM member WHERE MemberStatus != 'Removed' AND PrimaryEmail IS NOT NULL AND  Newsletter = 'electronic' ORDER BY LastName, FirstName")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for rows.Next() {
		var e string
		if err := rows.Scan(&e); err != nil {
			slog.Error(err.Error())
			// return err
			continue
		}
		r.Write([]string{"font@saint-luke.net", e, "USER", "MEMBER"})
	}
	r.Flush()
	return nil
}

func ReportSubscriber() ([]*Subscriber, error) {
	var subscribers []*Subscriber

	rows, err := db.Query("SELECT id FROM subscriber ORDER BY Name")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

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
}

// Returns a slice of IDs
func ActiveMemberIDs() ([]MemberID, error) {
	return reportMemberIDQuery("SELECT id FROM member WHERE MemberStatus != 'Removed' ORDER BY LastName, FirstName")
}

func reportMemberIDQuery(query string) ([]MemberID, error) {
	var id int
	list := make([]MemberID, 0, 500)

	rows, err := db.Query(query)
	if err != nil {
		slog.Error(err.Error())
		return list, err
	}

	for rows.Next() {
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
func ActiveSubscriberIDs() ([]SubscriberID, error) {
	return reportSubscriberIDQuery("SELECT id FROM subscriber WHERE DatePaid > DATE_SUB(CURRENT_DATE(), INTERVAL 366 DAY) ORDER BY Name")
}

func reportSubscriberIDQuery(query string) ([]SubscriberID, error) {
	var id int
	list := make([]SubscriberID, 0, 50)

	rows, err := db.Query(query)
	if err != nil {
		slog.Error(err.Error())
		return list, err
	}

	for rows.Next() {
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
	var members []*Member

	ids, err := ActiveMemberIDs()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, id := range ids {
		m, err := id.Get(true)
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
	r.Write([]string{"Name", "Address"})

	members, err := reportMemberIDQuery("SELECT id FROM member WHERE MemberStatus != 'Removed' AND DateReaffirmation > DATE_SUB(CURRENT_DATE(), INTERVAL 366 DAY) AND Doxology = 'mailed' ORDER BY LastName, FirstName")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, id := range members {
		m, err := id.Get(true)
		if err != nil {
			continue
		}
		member := []string{m.OSLName(), m.FormattedAddr}
		r.Write(member)
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
		subscriber := []string{s.Name, s.FormattedAddr}
		r.Write(subscriber)
	}
	r.Flush()
	return nil
}

// structured for Google Groups CSV upload
func DoxologyEmailed(w io.Writer) error {
	r := csv.NewWriter(w)
	r.Write([]string{"Group Email [Required]", "Member Email", "Member Type", "Member Role", "Name"})

	members, err := reportMemberIDQuery("SELECT id FROM member WHERE MemberStatus != 'Removed' AND Doxology = 'electronic' AND PrimaryEmail IS NOT NULL ORDER BY LastName, FirstName")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, id := range members {
		m, err := id.Get(true)
		if err != nil {
			continue
		}
		r.Write([]string{"doxology-distribution@saint-luke.net", m.PrimaryEmail, "USER", "MEMBER", m.OSLName()})
	}

	subscribers, err := reportSubscriberIDQuery("SELECT id FROM subscriber WHERE Doxology = 'electronic' AND PrimaryEmail IS NOT NULL DatePaid > DATE_SUB(CURRENT_DATE(), INTERVAL 730 DAY) ORDER BY Name")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, id := range subscribers {
		s, err := id.Get()
		if err != nil {
			continue
		}
		r.Write([]string{"doxology-distribution@saint-luke.net", s.PrimaryEmail, "USER", "MEMBER", s.Name + " : " + s.Attn})
	}
	r.Flush()
	return nil
}
