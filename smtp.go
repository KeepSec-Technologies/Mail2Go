package main

import (
	"fmt"
	"strings"

	"github.com/wneessen/go-mail"
)

// sendEmail constructs and sends an email with the specified HTML body and attachments.
func sendEmail(smtpServer string, smtpPort int, username string, password string, from string, to []string, subject, body, bodyFile string, attachments []string) error {
	// Create a new message
	m := mail.NewMsg()
	if err := m.From(from); err != nil {
		return err
	}

	// Set recipient(s)
	if err := m.To(to...); err != nil {
		return err
	}

	m.Subject(subject)

	// Set the body
	if bodyFile != "" {
		m.SetBodyString(mail.TypeTextHTML, body)
	} else {
		m.SetBodyString(mail.TypeTextPlain, body)
	}

	// Add attachments
	for _, attachment := range attachments {
		m.AttachFile(attachment)
	}

	// Create a new client
	c, err := mail.NewClient(smtpServer, mail.WithPort(smtpPort), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username), mail.WithPassword(password))
	if err != nil {
		return err
	}
	defer c.Close()

	// Send the email
	if err := c.DialAndSend(m); err != nil {
		fmt.Printf("Error sending email: %v", err)
		return err
	}

	fmt.Printf("\nEmail sent successfully to %s from %s\n", strings.Join(to, ", "), from)
	return nil
}
