package email

import (
	"context"
	"log/slog"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/matcornic/hermes/v2"
	"gopkg.in/gomail.v2"
)

func SendGeneric(ctx context.Context, ids []model.MemberID, subject string, message string) error {
	intros := strings.Split(message, "\n")
	messages := make([]*gomail.Message, 0, len(ids))

	for _, id := range ids {
		member, err := id.Get(ctx)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		if member.PrimaryEmail == "" {
			continue
		}

		e := hermes.Email{
			Body: hermes.Body{
				Name:   member.OSLShortName(),
				Intros: intros,
				Outros: []string{
					"Living the sacramental life",
				},
			},
		}

		body, err := hermesInstance.GenerateHTML(e)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		text, err := hermesInstance.GeneratePlainText(e)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		messages = append(messages, NewMessage(member.PrimaryEmail, subject, body, text))
	}

	if err := SendMany(messages...); err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
