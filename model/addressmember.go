package model

import (
	"github.com/Boostport/address"
)

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

// Address & Line 2 together
func (m *Member) formatPH() (address.Address, error) {
	return address.NewValid(
		address.WithCountry(m.Country),
		address.WithName(m.OSLName()),
		address.WithStreetAddress([]string{m.Address + ", " + m.AddressLine2}),
		address.WithLocality(m.City),
		address.WithAdministrativeArea(m.State),
		address.WithPostCode(m.PostalCode),
	)
}
