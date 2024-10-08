package email

import (
	"log/slog"

	"github.com/matcornic/hermes/v2"
)

func SendGiving(addr string, name string, amount string, description string) error {
	h, err := setup()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: name,
			Intros: []string{
				"Your donation has been recorded.",
			},
			Table: hermes.Table{
				Data: [][]hermes.Entry{
					{
						{Key: "Description", Value: description},
						{Key: "Amount", Value: amount},
					},
				},
			},
			Actions: []hermes.Action{
				{
					Instructions: "You can review your details in the online directory:",
					Button: hermes.Button{
						Text: "View your directory information",
						Link: "https://saint-luke.net/oo",
					},
				},
			},
			Outros: []string{
				"Living the sacramental life",
			},
		},
	}

	body, err := h.GenerateHTML(e)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	text, err := h.GeneratePlainText(e)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if err := send(addr, "OSL Donation Receipt", body, text); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
