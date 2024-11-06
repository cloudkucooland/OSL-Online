package email

import (
	"log/slog"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/matcornic/hermes/v2"
)

func SendGeneric(ids []model.MemberID, content string) error {
	h, err := setup()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	intros := strings.Split(content, "\n")

	for _, id := range ids {
		if err := sendGeneric(id, intros, h); err != nil {
			slog.Error(err.Error())
			// continue
		}
	}
	return nil
}

func sendGeneric(id model.MemberID, intros []string, h *hermes.Hermes) error {
	member, err := id.Get(true)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	if member.PrimaryEmail == "" {
		return nil
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

	if err := send(member.PrimaryEmail, "OSL Reaffirmation Reminder", body, text); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
