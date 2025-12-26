package email

import (
	"context"
	"log/slog"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/matcornic/hermes/v2"
)

func SendGeneric(ctx context.Context, ids []model.MemberID, subject string, message string) error {
	h, err := Setup()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	intros := strings.Split(message, "\n")

	for _, id := range ids {
		if err := sendGeneric(ctx, id, subject, intros, h); err != nil {
			slog.Error(err.Error())
			// continue
		}
	}
	return nil
}

func sendGeneric(ctx context.Context, id model.MemberID, subject string, intros []string, h *hermes.Hermes) error {
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

	if err := Send(member.PrimaryEmail, subject, body, text); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
