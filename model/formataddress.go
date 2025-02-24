package model

import (
	"log/slog"

	"github.com/Boostport/address"
)

// probably more trouble than it is worth to use an interface here
type addressFormatter interface {
	ISOCountry() string
	formatGB() (address.Address, error)
	formatPH() (address.Address, error)
	formatSG() (address.Address, error)
	formatMain() (address.Address, error)
}

// Make sure that addresses for mailing are properly formatted - handle both members and subscribers
func FormatAddress(m addressFormatter) (string, error) {
	var addr address.Address
	var err error

	switch m.ISOCountry() {
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
		slog.Error(err.Error(), "m", m)
		return "", err
	}

	postalStringFormatter := address.PostalLabelFormatter{
		Output:            address.StringOutputter{},
		OriginCountryCode: "US",
	}

	formatted := postalStringFormatter.Format(addr, "en")
	return formatted, nil
}
