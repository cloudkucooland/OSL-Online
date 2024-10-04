package email

import (
	"log/slog"

	"github.com/matcornic/hermes/v2"
)

func SendRegister(email string, password string) error {
	h, err := setup()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: email,
			Intros: []string{
				"OSL Online Directory Registration. Your password is " + password,
			},
			Actions: []hermes.Action{
				{
					Instructions: "To login, please click here:",
					Button: hermes.Button{
						Text: "Login",
						Link: "https://saint-luke.net/oo/#",
					},
				},
			},
			Outros: []string{
				"Living the sacramental life",
			},
		},
	}

	emailBody, err := h.GenerateHTML(e)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	emailText, err := h.GeneratePlainText(e)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if err := send(email, "OSL Direcgory register / password reset", emailBody, emailText); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
