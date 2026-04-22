package email

import (
	"crypto/tls"
	"net/mail"
	"os"
	"strconv"
	"time"

	"github.com/matcornic/hermes/v2"
	"gopkg.in/gomail.v2"
)

var (
	hermesInstance *hermes.Hermes
	smtpDialer     *gomail.Dialer
	fromAddress    mail.Address
)

func init() {
	email := os.Getenv("SENDER_EMAIL")
	if email == "" {
		email = "webmaster@saint-luke.net"
	}
	identity := os.Getenv("SENDER_IDENTITY")
	if identity == "" {
		identity = "OSL Webmaster"
	}
	fromAddress = mail.Address{Name: identity, Address: email}

	hermesInstance = &hermes.Hermes{
		Theme: new(hermes.Flat),
		Product: hermes.Product{
			Name:      "The Order of Saint Luke",
			Link:      "https://saint-luke.net/",
			Logo:      "https://saint-luke.net/static/logo.png",
			Copyright: "Copyright © " + strconv.Itoa(time.Now().Year()) + " The Order of Saint Luke. All rights reserved.",
		},
	}

	host := os.Getenv("SMTP_HOST")
	if host == "" {
		host = "localhost"
	}
	portStr := os.Getenv("SMTP_PORT")
	port, _ := strconv.Atoi(portStr)
	if port == 0 {
		port = 587
	}
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	smtpDialer = gomail.NewDialer(host, port, user, pass)

	skipVerify := true
	if sv := os.Getenv("SMTP_SKIP_VERIFY"); sv != "" {
		skipVerify, _ = strconv.ParseBool(sv)
	}
	smtpDialer.TLSConfig = &tls.Config{InsecureSkipVerify: skipVerify}
}

// Setup remains for backward compatibility but mostly returns the pre-initialized instance.
func Setup() (*hermes.Hermes, error) {
	return hermesInstance, nil
}

// Send sends a single email.
func Send(to string, subject string, htmlBody string, txtBody string) error {
	m := NewMessage(to, subject, htmlBody, txtBody)
	return smtpDialer.DialAndSend(m)
}

// NewMessage creates a new gomail.Message with the default From address.
func NewMessage(to string, subject string, htmlBody string, txtBody string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", fromAddress.String())
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/plain", txtBody)
	m.AddAlternative("text/html", htmlBody)
	return m
}

// SendMany sends multiple messages using a single SMTP connection.
func SendMany(messages ...*gomail.Message) error {
	if len(messages) == 0 {
		return nil
	}
	return smtpDialer.DialAndSend(messages...)
}

// GenerateAndSend generates HTML and text bodies from a hermes.Email and sends it.
func GenerateAndSend(to string, subject string, e hermes.Email) error {
	body, err := hermesInstance.GenerateHTML(e)
	if err != nil {
		return err
	}

	text, err := hermesInstance.GeneratePlainText(e)
	if err != nil {
		return err
	}

	return Send(to, subject, body, text)
}
