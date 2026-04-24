package email

import (
	"log/slog"

	"github.com/matcornic/hermes/v2"
)

func SendRegister(addr string, password string) error {
	e := hermes.Email{
		Body: hermes.Body{
			Title: "Registration & Password Reset",
			Name:  addr,
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

	if err := GenerateAndSend(addr, "OSL Directory Register / Password Reset requested", e); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
