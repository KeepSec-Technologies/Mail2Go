package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Long-form flags
	smtpServer string
	smtpPort   string
	username   string
	password   string

	fromEmail string
	toEmail   string

	subject string
	body    string

	attachmentsFiles string
	bodyFile         string

	// Short-form flags
	smtpServerShort string
	smtpPortShort   string
	usernameShort   string
	passwordShort   string

	fromEmailShort string
	toEmailShort   string

	subjectShort string
	bodyShort    string

	attachmentsFilesShort string
	bodyFileShort         string
)

func init() {
	// Long-form flags
	flag.StringVar(&smtpServer, "smtp-server", "", "SMTP server for sending emails")
	flag.StringVar(&smtpPort, "smtp-port", "587", "SMTP server port")
	flag.StringVar(&username, "smtp-username", "", "Username for SMTP authentication")
	flag.StringVar(&password, "smtp-password", "", "Password for SMTP authentication")

	flag.StringVar(&fromEmail, "from-email", "", "Email address to send from")
	flag.StringVar(&toEmail, "to-email", "", "Email addresses that will receive the email, comma-separated")

	flag.StringVar(&subject, "subject", "", "Subject of the email")
	flag.StringVar(&body, "body", "", "Body of the email")

	flag.StringVar(&attachmentsFiles, "attachments", "", "File paths for attachments, comma-separated")
	flag.StringVar(&bodyFile, "body-file", "", "File path for email body")

	// Short-form flags
	flag.StringVar(&smtpServerShort, "s", "", "SMTP server for sending emails (short)")
	flag.StringVar(&smtpPortShort, "p", "587", "SMTP server port (short)")
	flag.StringVar(&usernameShort, "u", "", "Username for SMTP authentication (short)")
	flag.StringVar(&passwordShort, "w", "", "Password for SMTP authentication (short)")

	flag.StringVar(&fromEmailShort, "f", "", "Email address to send from (short)")
	flag.StringVar(&toEmailShort, "t", "", "Email addresses that will receive the email, comma-separated (short)")

	flag.StringVar(&subjectShort, "h", "", "Subject of the email (short)")
	flag.StringVar(&bodyShort, "b", "", "Body of the email (short)")

	flag.StringVar(&attachmentsFilesShort, "af", "", "File paths for attachments, comma-separated (short)")
	flag.StringVar(&bodyFileShort, "bf", "", "File path for email body (short)")
}

func main() {
	flag.Parse()

	// Override long-form flags with short-form flags if set
	if smtpServerShort != "" {
		smtpServer = smtpServerShort
	}
	if smtpPortShort != "" {
		smtpPort = smtpPortShort
	}
	if usernameShort != "" {
		username = usernameShort
	}
	if passwordShort != "" {
		password = passwordShort
	}
	if fromEmailShort != "" {
		fromEmail = fromEmailShort
	}
	if toEmailShort != "" {
		toEmail = toEmailShort
	}
	if subjectShort != "" {
		subject = subjectShort
	}
	if bodyShort != "" {
		body = bodyShort
	}
	if attachmentsFilesShort != "" {
		attachmentsFiles = attachmentsFilesShort
	}
	if bodyFileShort != "" {
		bodyFile = bodyFileShort
	}

	// Check if required flags are missing
	if smtpServer == "" || username == "" || password == "" || fromEmail == "" || toEmail == "" || subject == "" {
		usage()
	}

	// Check if either direct input or file path is provided for body
	if body == "" && bodyFile == "" {
		fmt.Fprintln(os.Stderr, "Error: Subject and body are required, either directly or through file paths.")
		usage()
	}

	// Read body from files if provided
	if bodyFile != "" {
		content, err := os.ReadFile(bodyFile)
		if err != nil {
			log.Fatalf("\nError reading body file: %v\n", err)
		}
		body = string(content)
	}

	// Split attachment file paths
	var attachmentPaths []string
	if attachmentsFiles != "" {
		attachmentPaths = strings.Split(attachmentsFiles, ",")
	}

	// Split recipient email addresses
	var toEmails []string
	if toEmail != "" {
		toEmails = strings.Split(toEmail, ",")
	}

	if len(toEmails) == 0 {
		fmt.Fprintln(os.Stderr, "Error: At least one recipient email address is required.")
		usage()
	}

	sendEmail(fromEmail, toEmails, subject, body, attachmentPaths)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  -s, --smtp-server         SMTP server for sending emails")
	fmt.Fprintln(os.Stderr, "  -p, --smtp-port           SMTP server port (Default: 587)")
	fmt.Fprintln(os.Stderr, "  -u, --smtp-username       Username for SMTP authentication")
	fmt.Fprintln(os.Stderr, "  -w, --smtp-password       Password for SMTP authentication")
	fmt.Fprintln(os.Stderr, "  -f, --from-email          Email address to send notifications from")
	fmt.Fprintln(os.Stderr, "  -t, --to-email            Email addresses to send notifications to, comma-separated")
	fmt.Fprintln(os.Stderr, "  -h, --subject             Subject of the email")
	fmt.Fprintln(os.Stderr, "  -b, --body	            Body of the email")
	fmt.Fprintln(os.Stderr, "  -af, --attachments        File paths for attachments, comma-separated")
	fmt.Fprintln(os.Stderr, "  -bf, --body-file          File path for email body")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Ensure all required flags are provided.")
	os.Exit(1)
}

func getMimeType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}

// sendEmail constructs and sends an email with the specified HTML body and attachments.
func sendEmail(from string, to []string, subject, body string, attachments []string) error {
	// Set up authentication information
	auth := smtp.PlainAuth("", username, password, smtpServer)

	// Create the MIME multipart writer.
	var message bytes.Buffer
	writer := multipart.NewWriter(&message)

	// Set up the email headers.
	message.WriteString(fmt.Sprintf("From: %s\r\n", from))
	message.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary()))
	message.WriteString("\r\n")

	// Write the body part.
	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Type", "text/html; charset=utf-8") // Update this to "text/plain" if you are sending plain text
	partHeader.Set("Content-Transfer-Encoding", "quoted-printable")

	bodyPart, err := writer.CreatePart(partHeader)
	if err != nil {
		fmt.Printf("\nError creating body part: %v\n", err)
		return err
	}

	qw := quotedprintable.NewWriter(bodyPart)
	_, err = qw.Write([]byte(body))
	if err != nil {
		fmt.Printf("\nError writing body: %v\n", err)
		return err
	}
	qw.Close()

	// Attach files.
	for _, attachment := range attachments {
		file, err := os.Open(attachment)
		if err != nil {
			fmt.Printf("\nError opening attachment %s: %v\n", attachment, err)
			return err
		}
		defer file.Close()

		partHeader := textproto.MIMEHeader{}
		mimeType := getMimeType(attachment)
		partHeader.Set("Content-Type", mimeType)
		partHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(attachment)))
		partHeader.Set("Content-Transfer-Encoding", "base64")

		attachmentPart, err := writer.CreatePart(partHeader)
		if err != nil {
			fmt.Printf("\nError creating attachment part for %s: %v\n", attachment, err)
			return err
		}

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("\nError reading attachment %s: %v\n", attachment, err)
			return err
		}

		base64Writer := base64.NewEncoder(base64.StdEncoding, attachmentPart)
		_, err = base64Writer.Write(fileBytes)
		if err != nil {
			fmt.Printf("\nError encoding attachment %s: %v\n", attachment, err)
			return err
		}
		base64Writer.Close()
	}

	// Close the multipart writer to finalize the boundary.
	writer.Close()

	// Send the email.
	if err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, to, message.Bytes()); err != nil {
		fmt.Printf("Error sending email: %v", err)
		return err
	}

	fmt.Printf("\nEmail sent successfully to %s from %s\n", strings.Join(to, ", "), from)
	return nil
}
