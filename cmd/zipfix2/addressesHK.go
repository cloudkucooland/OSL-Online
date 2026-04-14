package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/cloudkucooland/OSL-Online/model"
)

// HKALSResponse matches the structure returned by ALS
type HKALSResponse struct {
	SuggestedAddress []struct {
		Address struct {
			PremisesAddress struct {
				EngPremisesAddress struct {
					BuildingName string `json:"BuildingName"`
					EngStreet    struct {
						StreetName     string `json:"StreetName"`
						BuildingNoFrom string `json:"BuildingNoFrom"`
					} `json:"EngStreet"`
					EngDistrict struct {
						DistrictName string `json:"DistrictName"`
					} `json:"EngDistrict"`
				} `json:"EngPremisesAddress"`
			} `json:"PremisesAddress"`
		} `json:"Address"`
		ValidationResult struct {
			Score float64 `json:"Score"`
		} `json:"ValidationResult"`
	} `json:"SuggestedAddress"`
}

func getaddressHK(ctx context.Context, member *model.Member) error {
	searchString := fmt.Sprintf("%s, %s, Hong Kong", member.Address, member.City)
	apiURL := "https://www.als.gov.hk/lookup?q=" + url.QueryEscape(searchString)

	req, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HK ALS status: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	var hkRes HKALSResponse
	if err := json.Unmarshal(body, &hkRes); err != nil {
		return err
	}

	if len(hkRes.SuggestedAddress) == 0 {
		slog.Info("no HK matches", "member", member.OSLName())
		return nil
	}

	topMatch := hkRes.SuggestedAddress[0]
	if topMatch.ValidationResult.Score < 70 {
		slog.Info("low confidence HK match", "member", member.OSLName(), "score", topMatch.ValidationResult.Score)
		return nil
	}

	addr := topMatch.Address.PremisesAddress.EngPremisesAddress

	// Format: "BuildingName, No StreetName"
	newStreet := fmt.Sprintf("%s %s", addr.EngStreet.BuildingNoFrom, addr.EngStreet.StreetName)
	if addr.BuildingName != "" {
		newStreet = fmt.Sprintf("%s, %s", addr.BuildingName, newStreet)
	}

	if member.Address != newStreet {
		slog.Info("correcting HK address", "old", member.Address, "new", newStreet)
		member.ID.SetMemberField(ctx, "Address", newStreet)
	}

	if member.City != addr.EngDistrict.DistrictName {
		member.ID.SetMemberField(ctx, "City", addr.EngDistrict.DistrictName)
	}

	return nil
}
