package email

import (
	"crypto/tls"
	"net/mail"
	"os"

	"github.com/matcornic/hermes/v2"
	"gopkg.in/gomail.v2"
)

var senderEmail string
var senderIdentity string

func setup() (*hermes.Hermes, error) {
	senderEmail = os.Getenv("SENDER_EMAIL")
	if senderEmail == "" {
		senderEmail = "wembaster@saint-luke.net"
	}
	senderIdentity = os.Getenv("SENDER_IDENTITY")
	if senderIdentity == "" {
		senderIdentity = "OSL Webmaster"
	}

	h := hermes.Hermes{
		Product: hermes.Product{
			Name: "The Order of Saint Luke",
			Link: "https://saint-luke.net/",
			Logo: "https://saint-luke.net/static/logo.png",
		},
	}
	return &h, nil
}

func send(to string, subject string, htmlBody string, txtBody string) error {
	from := mail.Address{
		Name:    senderIdentity,
		Address: senderEmail,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from.String())
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/plain", txtBody)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.Dialer{Host: "localhost", Port: 587}
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d.DialAndSend(m)
}
