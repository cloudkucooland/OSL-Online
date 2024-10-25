package model

import (
	"log/slog"

	"github.com/Boostport/address"
)

func (m *Member) FormatAddress() (string, error) {
	var addr address.Address
	var err error

	switch m.Country {
	case "GB":
		addr, err = m.formatGB()
	case "PH":
		addr, err = m.formatPH()
	case "SG":
		addr, err = m.formatSG()
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
func (m *Member) formatMain() (address.Address, error) {
	return address.NewValid(
		address.WithCountry(m.Country),
		address.WithName(m.OSLName()),
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
func (m *Member) formatGB() (address.Address, error) {
	return address.NewValid(
		address.WithCountry(m.Country),
		address.WithName(m.OSLName()),
		address.WithStreetAddress([]string{
			m.Address,
			m.AddressLine2,
		}),
		address.WithLocality(m.City),
		address.WithPostCode(m.PostalCode),
	)
}

// no AdministrativeArea or Locality
func (m *Member) formatSG() (address.Address, error) {
	return address.NewValid(
		address.WithCountry(m.Country),
		address.WithName(m.OSLName()),
		address.WithStreetAddress([]string{
			m.Address,
			m.AddressLine2,
		}),
		address.WithPostCode(m.PostalCode),
	)
}

func (m *Member) formatPH() (address.Address, error) {
	return address.NewValid(
		address.WithCountry(m.Country),
		address.WithName(m.OSLName()),
		address.WithStreetAddress([]string{m.Address + ", " + m.AddressLine2}),
		address.WithAdministrativeArea(m.State),
		address.WithLocality(m.City),
		address.WithPostCode(m.PostalCode),
	)
}
