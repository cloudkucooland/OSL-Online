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
		return m.formatUS()
	case "Phillipines", "PHILIPPINES", "Philipines", "Philippines":
		m.Country = "PH"
		m.Store()
		return "", nil // for now
	case "UNITED KINGDOM", "United Kingdom":
		m.Country = "GB"
		m.Store()
		return m.formatGB()
	case "CANADA":
		m.Country = "CA"
		m.Store()
		return m.formatUS()
	case "Hong Kong":
		m.Country = "HK"
		m.Store()
		return m.formatUS()
	case "SINGAPORE", "Singapore":
		m.Country = "SG"
		m.Store()
		return "", nil // for now
	case "GB":
		return m.formatGB()
	case "PH", "SG":
		return "", nil // for now
	default: // assume US/CA/HK format
		return m.formatUS()
	}
	return m.formatUS()
}

func (m *Member) formatUS() (string, error) {
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
