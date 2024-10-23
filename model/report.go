package model

import (
	"encoding/csv"

	"io"
	"log/slog"
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

func ReportNotRenewed() ([]*Member, error) {
	return reportMemberQuery("SELECT id FROM member WHERE MemberStatus = 'Annual Vows' AND DateReaffirmation < DATE_SUB(CURRENT_DATE(), INTERVAL 365 DAY) ORDER BY DateReaffirmation")
}

func ReportExpired() ([]*Member, error) {
	return reportMemberQuery("SELECT id FROM member WHERE MemberStatus = 'Annual Vows' AND DateReaffirmation < DATE_SUB(CURRENT_DATE(), INTERVAL 730 DAY) ORDER BY DateReaffirmation")
}

func ReportAnnual() ([]*Member, error) {
	return reportMemberQuery("SELECT id FROM member WHERE MemberStatus = 'Annual Vows' ORDER BY LastName, FirstName")
}

func ReportLife() ([]*Member, error) {
	return reportMemberQuery("SELECT id FROM member WHERE MemberStatus = 'Life Vows' ORDER BY LastName, FirstName")
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
		addr, err := m.FormatAddress()
		if err != nil {
			continue
		}
		member := []string{m.OSLName(), addr}
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
		addr, _ := s.FormatAddress()
		if err != nil {
			continue
		}
		subscriber := []string{s.Name, addr}
		r.Write(subscriber)
	}
	r.Flush()
	return nil
}

func DoxologyEmailed(w io.Writer) error {
	r := csv.NewWriter(w)
	r.Write([]string{"Name", "Email"})

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
		member := []string{m.OSLName(), m.PrimaryEmail}
		r.Write(member)
	}

	subscribers, err := reportSubscriberIDQuery("SELECT id FROM subscriber WHERE Doxology = 'electronic' AND PrimaryEmail IS NOT NULL ORDER BY Name")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, id := range subscribers {
		s, err := id.Get()
		if err != nil {
			continue
		}
		subscriber := []string{s.Name, s.PrimaryEmail}
		r.Write(subscriber)
	}
	r.Flush()
	return nil
}
