package model

import (
	"log/slog"

	"github.com/Boostport/address"
)

func (m *Subscriber) FormatAddress() (string, error) {
	var addr address.Address
	var err error

	switch m.Country {
	case "GB":
		addr, err = m.formatGB()
	case "PH":
		addr, err = m.formatPH()
	case "SG":
		addr, err = m.formatSG() // for now
	default: // assume US/CA/HK format
		addr, err = m.formatMain()
	}

	if err != nil {
		slog.Error(err.Error(), m)
		return "", err
	}

	postalStringFormatter := address.PostalLabelFormatter{
		Output:            address.StringOutputter{},
		OriginCountryCode: "US",
	}

	formatted := postalStringFormatter.Format(addr, "en")
	return formatted, nil
}

// this does the "RIGHT THING"TM for most countries, even HK which goes MSB vs. US's LSB
func (m *Subscriber) formatMain() (address.Address, error) {
	return address.NewValid(
		address.WithCountry(m.Country),
		address.WithName(m.Attn),
		address.WithOrganization(m.Name),
		address.WithStreetAddress([]string{
			m.Address,
			m.AddressLine2,
		}),
		address.WithLocality(m.City),
		address.WithAdministrativeArea(m.State),
		address.WithPostCode(m.PostalCode),
	)
}

// no AdministrativeArea
func (m *Subscriber) formatGB() (address.Address, error) {
	return address.NewValid(
		address.WithCountry(m.Country),
		address.WithName(m.Attn),
		address.WithOrganization(m.Name),
		address.WithStreetAddress([]string{
			m.Address,
			m.AddressLine2,
		}),
		address.WithLocality(m.City),
		address.WithPostCode(m.PostalCode),
	)
}

// no AdministrativeArea or Locality
func (m *Subscriber) formatSG() (address.Address, error) {
	return address.NewValid(
		address.WithCountry(m.Country),
		address.WithOrganization(m.Name),
		address.WithName(m.Attn),
		address.WithStreetAddress([]string{
			m.Address,
			m.AddressLine2,
		}),
		address.WithPostCode(m.PostalCode),
	)
}

func (m *Subscriber) formatPH() (address.Address, error) {
	return address.NewValid(
		address.WithCountry(m.Country),
		address.WithOrganization(m.Name),
		address.WithName(m.Attn),
		address.WithStreetAddress([]string{m.Address + ", " + m.AddressLine2}),
		address.WithAdministrativeArea(m.State),
		address.WithLocality(m.City),
		address.WithPostCode(m.PostalCode),
	)
}
