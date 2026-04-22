package email

import (
	"context"
	"log/slog"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/matcornic/hermes/v2"
)

func SendGiving(ctx context.Context, id model.MemberID, amount string, description string) error {
	member, err := id.Get(ctx)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	if member.PrimaryEmail == "" {
		return nil
	}

	e := hermes.Email{
		Body: hermes.Body{
			Title: "Donation Receipt",
			Name:  member.OSLName(),
			Intros: []string{
				"Your donation has been recorded. Thank you for your support of the Order.",
			},
			Table: hermes.Table{
				Data: [][]hermes.Entry{
					{
						{Key: "Description", Value: description},
						{Key: "Amount", Value: "$" + amount},
					},
				},
			},
			Actions: []hermes.Action{
				{
					Instructions: "You can review your details in the online directory:",
					Button: hermes.Button{
						Text: "View your directory information",
						Link: "https://saint-luke.net/oo/#/me",
					},
				},
			},
			Outros: []string{
				"Living the sacramental life",
			},
		},
	}

	if err := GenerateAndSend(member.PrimaryEmail, "OSL Donation Receipt", e); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
