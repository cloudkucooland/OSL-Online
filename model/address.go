package model

import (
	"log/slog"

	"github.com/Boostport/address"
)

func (m *Member) FormatAddress() (string, error) {
	switch m.Country {
	case "USA", "United States", "":
		m.Country = "US"
		m.Store()
		return m.formatMain()
	case "Phillipines", "PHILIPPINES", "Philipines", "Philippines":
		m.Country = "PH"
		m.Store()
		return m.formatPH()
	case "UNITED KINGDOM", "United Kingdom":
		m.Country = "GB"
		m.Store()
		return m.formatGB()
	case "CANADA":
		m.Country = "CA"
		m.Store()
		return m.formatMain()
	case "Hong Kong":
		m.Country = "HK"
		m.Store()
		return m.formatMain()
	case "SINGAPORE", "Singapore":
		m.Country = "SG"
		m.Store()
		return m.formatSG()
	case "GB":
		return m.formatGB()
	case "PH":
		return m.formatPH()
	case "SG":
		return m.formatSG() // for now
	default: // assume US/CA/HK format
		return m.formatMain()
	}
	// not-reached
	return m.formatMain()
}

// this does the "RIGHT THING"TM for most countries, even HK which goes MSB vs. US's LSB
func (m *Member) formatMain() (string, error) {
	addr, err := address.NewValid(
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
	if err != nil {
		slog.Error(err.Error(), "data", m)
		return "", err
	}

	postalStringFormatter := address.PostalLabelFormatter{
		Output:            address.StringOutputter{},
		OriginCountryCode: "US",
	}

	formatted := postalStringFormatter.Format(addr, "en")
	return formatted, nil
}

// no AdministrativeArea
func (m *Member) formatGB() (string, error) {
	addr, err := address.NewValid(
		address.WithCountry(m.Country),
		address.WithName(m.OSLName()),
		address.WithStreetAddress([]string{
			m.Address,
			m.AddressLine2,
		}),
		address.WithLocality(m.City),
		address.WithPostCode(m.PostalCode),
	)
	if err != nil {
		slog.Error(err.Error(), "data", m)
		return "", err
	}

	postalStringFormatter := address.PostalLabelFormatter{
		Output:            address.StringOutputter{},
		OriginCountryCode: "US",
	}

	formatted := postalStringFormatter.Format(addr, "en")
	return formatted, nil
}

// no AdministrativeArea or Locality
func (m *Member) formatSG() (string, error) {
	addr, err := address.NewValid(
		address.WithCountry(m.Country),
		address.WithName(m.OSLName()),
		address.WithStreetAddress([]string{
			m.Address,
			m.AddressLine2,
		}),
		address.WithPostCode(m.PostalCode),
	)
	if err != nil {
		slog.Error(err.Error(), "data", m)
		return "", err
	}

	postalStringFormatter := address.PostalLabelFormatter{
		Output:            address.StringOutputter{},
		OriginCountryCode: "US",
	}

	formatted := postalStringFormatter.Format(addr, "en")
	return formatted, nil
}

func (m *Member) formatPH() (string, error) {
	addr, err := address.NewValid(
		address.WithCountry(m.Country),
		address.WithName(m.OSLName()),
		address.WithStreetAddress([]string{m.Address + ", " + m.AddressLine2}),
		address.WithAdministrativeArea(m.State),
		address.WithLocality(m.City),
		address.WithPostCode(m.PostalCode),
	)
	if err != nil {
		slog.Error(err.Error(), "data", m)
		return "", err
	}

	postalStringFormatter := address.PostalLabelFormatter{
		Output:            address.StringOutputter{},
		OriginCountryCode: "US",
	}

	formatted := postalStringFormatter.Format(addr, "en")
	return formatted, nil
}
