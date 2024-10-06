package model

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/Boostport/address"
	"github.com/hashicorp/go-multierror"
)

func (m *Member) ValidateAddress() (string, error) {
	name := fmt.Sprintf("%s %s", m.FirstName, m.LastName)

	addr, err := address.NewValid(
		address.WithCountry(m.Country),
		address.WithName(name),
		address.WithStreetAddress([]string{
			m.Address,
			m.AddressLine2,
		}),
		address.WithLocality(m.City),
		address.WithAdministrativeArea(m.State),
		address.WithPostCode(m.PostalCode),
	)

	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	if err != nil {
		// If there was an error and you want to find out which validations failed,
		// type switch the nested error as a *multierror.Error to access the list of errors
		if merr, ok := errors.Unwrap(err).(*multierror.Error); ok {
			for _, subErr := range merr.Errors {
				if subErr == address.ErrInvalidCountryCode {
					slog.Error(subErr.Error())
				}
			}
		}
		return "", err
	}

	postalStringFormatter := address.PostalLabelFormatter{
		Output:            address.StringOutputter{},
		OriginCountryCode: "US",
	}

	formatted := postalStringFormatter.Format(addr, "en")
	fmt.Println(formatted)
	return formatted, nil
}
