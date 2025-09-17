package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// oauthurl = "https://apis-tem.usps.com/oauth2/v3"
// https://apis.usps.com/oauth2/v3

type Bearer struct {
	AccessToken     string `json:"access_token"`
	TokenType       string `json:"token_type"`
	IssuedAt        int    `json:"issued_at"`
	ExpiresIn       int    `json:"exires_in"`
	Status          string `json:"status"`
	Scope           string `json:"scope"`
	Issuer          string `json:"issuer"`
	ClientID        string `json:"client_id"`
	ApplicationName string `json:"application_name"`
	APIProducts     string `json:"api_products"`
	PublicKey       string `json:"public_key"`
}

// roll our own since USPS requires POST and that's not-standard
func getauth(ctx context.Context) string {
	clientID := os.Getenv("ZIPFIX_CLIENTID")
	clientSecret := os.Getenv("ZIPFIX_SECRET")

	if clientID == "" || clientSecret == "" {
		panic("ZIPFIX_CLIENTID and ZIPFIX_SECRET must be set")
	}

	j := fmt.Sprintf("{\"client_id\": \"%s\", \"client_secret\":\"%s\", \"grant_type\":\"client_credentials\"}", clientID, clientSecret)
	authURL := "https://apis.usps.com/oauth2/v3/token"

	client := &http.Client{}

	req, err := http.NewRequest("POST", authURL, strings.NewReader(j))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	bearer := Bearer{}
	err = json.Unmarshal(body, &bearer)
	if err != nil {
		panic(err)
	}
	return string(bearer.AccessToken)
}
