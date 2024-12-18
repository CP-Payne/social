package mailer

import (
	"bytes"
	"fmt"
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

	fromEmail := m.fromEmail
	toEmail := email

	// If dev environment, send email to fromEmail instead of toEmail
	// Email will be sent from mailtrap
	if isSandbox {
		fromEmail = mailTrapDemoMail
		toEmail = m.fromEmail
	}

	message := gomail.NewMessage()
	message.SetHeader("From", fromEmail)
	message.SetHeader("To", toEmail)

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

	// Maybe remove the retry into a separate function in the handler. Using custom logger would then be easier
	var retryErr error
	for i := 0; i < maxRetries; i++ {
		retryErr = m.client.DialAndSend(message)
		if retryErr != nil {
			// log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetries)
			// log.Printf("Error: %v", err.Error())

			// Exponential Backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		return nil
	}

	return fmt.Errorf("failed to send email after %d attempt, error: %v", maxRetries, retryErr)
}
