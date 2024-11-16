package model

import (
	"github.com/nyaruka/phonenumbers"
)

func FormatPhoneNumber(number string, country string) (string, error) {
	if number == "" {
		return "", nil
	}

	p, err := phonenumbers.Parse(number, country)
	if err != nil {
		return "", err
	}

	mode := phonenumbers.INTERNATIONAL
	if country == "US" || country == "CA" {
		mode = phonenumbers.NATIONAL
	}
	formatted := phonenumbers.Format(p, mode)
	return formatted, nil
}
