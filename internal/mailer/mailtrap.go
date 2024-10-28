package mailer

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"

	gomail "gopkg.in/mail.v2"
)

var (
	liveMailTrapHost = "live.smtp.mailtrap.io"
	mailTrapPort     = 587
	mailTrapDemoMail = "hello@demomailtrap.com"
)

type MailTrapMailer struct {
	fromEmail string
	client    gomail.Dialer
}

func NewMailTrap(username, password, fromEmail string) *MailTrapMailer {
	dialer := gomail.NewDialer(liveMailTrapHost, mailTrapPort, username, password)
	return &MailTrapMailer{
		fromEmail: fromEmail,
		client:    *dialer,
	}
}

func (m *MailTrapMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	message := gomail.NewMessage()

	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)

	// If dev environment, send email to fromEmail instead of toEmail
	// Email will be sent from mailtrap
	if isSandbox {
		message.SetHeader("From", mailTrapDemoMail)
		message.SetHeader("To", m.fromEmail)
	}
	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	message.SetHeader("Subject", subject.String())
	message.SetBody("text/html", body.String())

	for i := 0; i < maxRetries; i++ {
		err = m.client.DialAndSend(message)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetries)
			log.Printf("Error: %v", err.Error())

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		log.Printf("Email sent")
		return nil
	}

	return fmt.Errorf("failed to send emaila fter %d attempts", maxRetries)
}
