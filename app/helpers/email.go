package helpers

import (
	"gopkg.in/gomail.v2"
	"log"
	"os"
)

const CONFIG_SMTP_PORT = 587

func SendMail(email string, body string, subject string) {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("CONFIG_EMAIL"))
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		os.Getenv("CONFIG_SMTP_HOST"),
		CONFIG_SMTP_PORT,
		os.Getenv("CONFIG_EMAIL"),
		os.Getenv("CONFIG_PASSWORD"),
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Mail sent!")
}
