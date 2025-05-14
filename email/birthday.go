package email

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/matcornic/hermes/v2"
)

type BirthdayEmailEntry struct {
	ID   int
	Name string
}

func SendBirthdayMail(members []*BirthdayEmailEntry, month time.Month, day int) error {
	if len(members) == 0 {
		slog.Info("no birthdays today")
		return nil
	}

	h, err := Setup()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	data := [][]hermes.Entry{
		{},
	}
	for _, m := range members {
		entry := hermes.Entry{
			Key:   m.Name,
			Value: fmt.Sprintf("https://saint-luke.net/oo/#/member/%d", m.ID),
		}
		data[0] = append(data[0], entry)
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: "Siblings",
			Intros: []string{
				"Today's OSL Birthdays",
			},
			Table: hermes.Table{
				Data: data,
			},
			Actions: []hermes.Action{
				{
					Instructions: "Member details in the directory:",
					Button: hermes.Button{
						Text: "OSL Directory",
						Link: "https://saint-luke.net/oo/#",
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

	if err := Send("birthdays@saint-luke.net", "Today's OSL Birthdays", body, text); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
