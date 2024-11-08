package email

import (
	"log/slog"

	"github.com/matcornic/hermes/v2"
)

func SendRegister(addr string, password string) error {
	h, err := Setup()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: addr,
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

	if err := Send(addr, "OSL Directory Register / Password Reset requested", body, text); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
