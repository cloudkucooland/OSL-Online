package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"
)

type OneMapTokenResponse struct {
	Token string `json:"access_token"`
}

type OneMapSearchResponse struct {
	Found   int `json:"found"`
	Results []struct {
		Address    string `json:"ADDRESS"`
		RoadName   string `json:"ROAD_NAME"`
		Block      string `json:"BLK_NO"`
		Building   string `json:"BUILDING"`
		PostalCode string `json:"POSTAL"`
	} `json:"results"`
}

func getOneMapToken() (string, error) {
	email := os.Getenv("ONEMAP_EMAIL")
	password := os.Getenv("ONEMAP_PASSWORD")

	payload := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password)
	resp, err := http.Post("https://www.onemap.gov.sg/api/auth/post/getToken", "application/json", strings.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tr OneMapTokenResponse
	json.NewDecoder(resp.Body).Decode(&tr)
	return tr.Token, nil
}

func getaddressSG(ctx context.Context, member *model.Member, token string) error {
	// OneMap works best if you just search the postal code
	query := member.PostalCode
	apiURL := fmt.Sprintf("https://www.onemap.gov.sg/api/common/elastic/search?searchVal=%s&returnGeom=N&getAddrDetails=Y", url.QueryEscape(query))

	req, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var sr OneMapSearchResponse
	json.NewDecoder(resp.Body).Decode(&sr)

	if sr.Found > 0 {
		res := sr.Results[0]
		// Format: "BLK 123 ROAD NAME, #UNIT-NO (if available) BUILDING"
		// OneMap doesn't know the Unit #, so we keep what the user had if possible
		newAddr := fmt.Sprintf("BLK %s %s", res.Block, res.RoadName)
		if res.Building != "NIL" && res.Building != "" {
			newAddr = fmt.Sprintf("%s, %s", newAddr, res.Building)
		}

		if member.Address != newAddr {
			member.ID.SetMemberField(ctx, "Address", newAddr)
		}
		// SG is small, City is always Singapore
		if member.City != "Singapore" {
			member.ID.SetMemberField(ctx, "City", "Singapore")
		}
	}
	return nil
}
