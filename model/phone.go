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

	pf := ""
	if country == "US" || country == "CA" {
		pf = phonenumbers.Format(p, phonenumbers.NATIONAL)
	} else {
		pf = phonenumbers.Format(p, phonenumbers.INTERNATIONAL)
	}
	return pf, nil
}
