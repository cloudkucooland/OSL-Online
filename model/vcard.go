package model

import (
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/emersion/go-vcard"
)

func (m *Member) WriteVCard(w io.Writer) error {
	card := make(vcard.Card)

	card.SetValue(vcard.FieldFormattedName, m.OSLName())
	card.SetValue(vcard.FieldNickname, m.OSLShortName())
	card.SetValue(vcard.FieldBirthday, m.BirthDate.Format(time.DateOnly))
	card.SetValue(vcard.FieldTelephone, m.PrimaryPhone)
	card.SetValue(vcard.FieldEmail, m.PrimaryEmail)

	name := vcard.Name{
		FamilyName:      m.LastName,
		GivenName:       m.FirstName,
		AdditionalName:  m.MiddleName,
		HonorificPrefix: m.Title,
		HonorificSuffix: "OSL",
	}
	card.AddName(&name)

	addr := vcard.Address{
		ExtendedAddress: m.AddressLine2,
		StreetAddress:   m.Address,
		Locality:        m.City,
		Region:          m.State,
		PostalCode:      m.PostalCode,
		Country:         m.Country,
	}
	card.AddAddress(&addr)

	vcard.ToV4(card)
	e := vcard.NewEncoder(w)
	if err := e.Encode(card); err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
