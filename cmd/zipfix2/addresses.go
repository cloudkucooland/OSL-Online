package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/cloudkucooland/OSL-Online/model"
)

const URL = "https://apis.usps.com/addresses/v3/address"

func getaddress(ctx context.Context, member *model.Member, bearer string) {
	if member.Address == "" || member.State == "" || member.City == "" {
		slog.Info("getaddress", "not enough data", member.OSLName())
		return
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", `Bearer `+bearer)

	values := req.URL.Query()
	values.Add("streetAddress", member.Address)
	if member.AddressLine2 != "" {
		values.Add("secondaryAddress", member.AddressLine2)
	}
	values.Add("city", member.City)
	values.Add("state", member.State)

	if member.PostalCode != "" {
		values.Add("ZIPCode", member.PostalCode[:5])

		if len(member.PostalCode) == 10 {
			values.Add("ZIPPlus4", member.PostalCode[6:9])
		}
	}

	req.URL.RawQuery = values.Encode()

RETRY:
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	switch resp.StatusCode {
	case 200:
		// nothing
	case 429:
		slog.Info("going to fast, backing off for 5 minutes")
		time.Sleep(300 * time.Second)
		goto RETRY
	default:
		slog.Info("bad status", "status", resp.Status)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	ar := AddressResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		panic(err)
	}

	if len(ar.Warnings) != 0 {
		slog.Error("Warnings found", "warnings", ar.Warnings, "member", member.OSLName())
		slog.Info("ar", "ar", ar)
		return
	}

	if len(ar.Matches) == 0 {
		slog.Info("no matches", "member", member.OSLName())
		return
	}

	for i := range ar.Matches {
		switch ar.Matches[i].Code {
		case "31", "":
			break
			// good result
		case "32":
			// need apt or suite number
			slog.Error("Need more information", "code", ar.Matches[0].Code, "member", member.OSLName(), "info", ar.Matches[0].Text)
			// continue, we can still fix the zip
			break
		default:
			slog.Error("Unknown Matches code", "code", ar.Matches[0].Code, "member", member.OSLName(), "info", ar.Matches[0].Text)
			slog.Info("ar", "ar", ar)
			continue
		}
	}

	for i := range ar.Corrections {
		switch ar.Corrections[1].Code { // need to get a list of codes...
		case "":
			slog.Info("corrections with empty code", i, ar.Corrections[i])
		default:
			slog.Info("unknown corrections", i, ar.Corrections[i])
		}
		// continue, we can still fix the zip
	}

	if ar.Additional.Vacant == "Y" {
		slog.Info("Address is marked as Vacant", "member", member.OSLName())
		// continue, we can still fix the zip
	}

	if ar.Address.ZIPCode != "" && ar.Address.ZIPPlus4 != "" {
		z := fmt.Sprintf("%s-%s", ar.Address.ZIPCode, ar.Address.ZIPPlus4)
		if member.PostalCode != z {
			slog.Info("updating", "name", member.OSLName(), "current ZIP", member.PostalCode, "new ZIP", z)
		}
		member.PostalCode = z
		if err := member.Store(); err != nil {
			slog.Error("unable to store", "error", err.Error(), "member", member.OSLName())
		}
	}
}

type AddressResponse struct {
	Firm        string        `json:"firm"`
	Address     Address       `json:"address"`
	Additional  Additional    `json:"additionalInfo"`
	Corrections []Corrections `json:"corrections"`
	Matches     []Matches     `json:"matches"`
	Warnings    []string      `json:"warnings"`
}

type Address struct {
	StreetAddress    string `json:"streetAddress"`
	StreetAbbr       string `json:"streeAddressAbbreviation"`
	SecondaryAddress string `json:"secondaryAddress"`
	CityAbbreviation string `json:"cityAbbreviation"`
	City             string `json:"city"`
	State            string `json:"state"`
	ZIPCode          string
	ZIPPlus4         string
	Urbanization     string `json:"urbanization"`
}

type Additional struct {
	DeliveryPoint   string `json:"deliveryPoint"`
	CarrierRoute    string `json:"carrierRoute"`
	DPVConfirmation string
	DPVCRMA         string
	Business        string `json:"business"`
	CentralDelivery string `json:"centralDeliveryPoint"`
	Vacant          string `json:"vacant"`
}

type Corrections struct {
	Code string `json:"code"`
	Text string `json:"text"`
}

type Matches struct {
	Code string `json:"code"`
	Text string `json:"text"`
}
