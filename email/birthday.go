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

	data := make([][]hermes.Entry, 0, len(members))
	for _, m := range members {
		data = append(data, []hermes.Entry{
			{Key: "Member", Value: m.Name},
			{Key: "Directory Link", Value: fmt.Sprintf("https://saint-luke.net/oo/#/member/%d", m.ID)},
		})
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: "Siblings",
			Intros: []string{
				fmt.Sprintf("Today's OSL Birthdays for %s %d", month.String(), day),
			},
			Table: hermes.Table{
				Data: data,
				Columns: hermes.Columns{
					CustomWidth: map[string]string{
						"Member":         "40%",
						"Directory Link": "60%",
					},
					CustomAlignment: map[string]string{
						"Directory Link": "left",
					},
				},
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

	if err := GenerateAndSend("birthdays@saint-luke.net", "Today's OSL Birthdays", e); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
