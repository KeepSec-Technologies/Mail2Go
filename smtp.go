package main

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/wneessen/go-mail"
)

// sendEmail constructs and sends an email with the specified HTML body and attachments.
func sendEmail(smtpServer string, smtpPort int, username string, password string, from string, to []string, subject, body, bodyFile string, attachments []string, tlsMode string) error {
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

	clientOptions := []mail.Option{
		mail.WithPort(smtpPort),
	}

	// Define client options
	if username != "" || password != "" {
		clientOptions = append(
			clientOptions,
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(username),
			mail.WithPassword(password),
		)
	}

	// Conditionally add TLS options based on tlsMode
	switch tlsMode {
	case "none":
		clientOptions = append(clientOptions, mail.WithTLSPolicy(mail.NoTLS))
	case "tls-skip":
		tlsSkipConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpServer,
		}
		clientOptions = append(clientOptions, mail.WithTLSConfig(tlsSkipConfig))
	case "tls":
		clientOptions = append(clientOptions, mail.WithTLSPolicy(mail.TLSMandatory))
	}
	// Create a new client using the options
	c, err := mail.NewClient(smtpServer, clientOptions...)
	if err != nil {
		fmt.Printf("Failed to create SMTP client: %v", err)
	}
	if c == nil {
		fmt.Printf("SMTP client is nil")
	}

	// Send the email
	if err := c.DialAndSend(m); err != nil {
		fmt.Printf("Error sending email: %v", err)
		return err
	}

	fmt.Printf("\nEmail sent successfully to %s from %s\n", strings.Join(to, ", "), from)
	return nil
}
