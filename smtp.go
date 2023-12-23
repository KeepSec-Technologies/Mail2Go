package main

import (
	"fmt"
	"strings"

	"gopkg.in/mail.v2"
)

// sendEmail constructs and sends an email with the specified HTML body and attachments.
func sendEmail(smtpServer string, smtpPort int, username string, password string, from string, to []string, subject, body, bodyFile string, attachments []string) error {
	m := mail.NewMessage()

	// Set the main email parts
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	// Check if the body is provided via a file
	if bodyFile != "" {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	// Add attachments
	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	// Set up SMTP information
	d := mail.NewDialer(smtpServer, smtpPort, username, password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Error sending email: %v", err)
		return err
	}

	fmt.Printf("\nEmail sent successfully to %s from %s\n", strings.Join(to, ", "), from)
	return nil
}
