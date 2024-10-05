package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"

	"github.com/aureum/usps-go"
)

func main() {
	ctx := context.Background()
	var u usps.USPS
	u.Username = os.Getenv("ZIPFIX_USER")
	if u.Username == "" {
		panic("ZIPFIX_USER not set")
	}
	u.Password = os.Getenv("ZIPFIX_PASS")
	if u.Password == "" {
		panic("ZIPFIX_PASS not set")
	}
	// u.Production = true

	dbpath := os.Getenv("OO_DB")
	if dbpath == "" {
		panic("OO_DB enviornment var not set. e.g. oo:password@unix(/var/lib/mysql/mysql.sock)/oo")
	}

	if err := model.Connect(ctx, dbpath); err != nil {
		slog.Error("startup", "message", "Error connecting to database", "error", err.Error())
		panic(err)
	}
	ids, err := model.ActiveMembers()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	var a usps.Address
	for _, i := range ids {
		m, err := model.GetMember(i, true)
		if err != nil {
			slog.Error(err.Error())
			panic(err)
		}
		fmt.Printf("%s\n", m.OSLName())

		if m.Country == "United States" || m.Country == "USA" || m.Country == "US" {
			m.Country = "US"
		} else {
			fmt.Printf("Skipping non-us: %s\n", m.Country)
			continue
		}

		// a.Address1 = m.Address
		// a.Address2 = m.AddressLine2
		a.Address2 = m.Address + " " + m.AddressLine2
		a.City = m.City
		a.State = m.State
		a.Zip5 = zip4(m.PostalCode)
		fmt.Printf("%+v\n", a)
		lookup := u.ZipCodeLookup(a)
		fmt.Printf("%+v\n", lookup)

		// verify := u.AddressVerification(a)
		// fmt.Printf("%+v\n", verify)

		newzip := fmt.Sprintf("%s-%s", lookup.Address.Zip5, lookup.Address.Zip4)
		if lookup.Address.Zip5 != "" && newzip != m.PostalCode {
			fmt.Printf("Old: %s\tNew: %s", m.PostalCode, newzip)
			// model.SetField(i, "PostalCode", newzip)
		}
	}
}

func zip4(z string) string {
	if z == "" {
		return ""
	}
	zs := strings.Split(z, "-")
	return zs[0]
}
