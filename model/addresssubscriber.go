package model

import (
	"github.com/Boostport/address"
)

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
