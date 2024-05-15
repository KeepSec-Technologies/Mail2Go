package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	// Long-form flags
	smtpServer string
	smtpPort   int
	username   string
	password   string

	tlsMode string

	configFile string

	fromEmail string
	toEmail   string

	subject string
	body    string

	attachmentsFiles string
	bodyFile         string

	// Short-form flags
	smtpServerShort string
	smtpPortShort   int
	usernameShort   string
	passwordShort   string

	tlsModeShort string

	configFileShort string

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
	flag.IntVar(&smtpPort, "smtp-port", 587, "SMTP server port")
	flag.StringVar(&username, "smtp-username", "", "Username for SMTP authentication")
	flag.StringVar(&password, "smtp-password", "", "Password for SMTP authentication")

	flag.StringVar(&tlsMode, "tls-mode", "", "TLS mode (none, tls-skip, tls)")

	flag.StringVar(&configFile, "config", "", "Path to the SMTP config file")

	flag.StringVar(&fromEmail, "from-email", "", "Email address to send from")
	flag.StringVar(&toEmail, "to-email", "", "Email addresses that will receive the email, comma-separated")

	flag.StringVar(&subject, "subject", "", "Subject of the email")
	flag.StringVar(&body, "body", "", "Body of the email")

	flag.StringVar(&attachmentsFiles, "attachments", "", "File paths for attachments, comma-separated")
	flag.StringVar(&bodyFile, "body-file", "", "File path for email body")

	// Short-form flags
	flag.StringVar(&smtpServerShort, "s", "", "SMTP server for sending emails (short)")
	flag.IntVar(&smtpPortShort, "p", 587, "SMTP server port (short)")
	flag.StringVar(&usernameShort, "u", "", "Username for SMTP authentication (short)")
	flag.StringVar(&passwordShort, "w", "", "Password for SMTP authentication (short)")

	flag.StringVar(&tlsModeShort, "l", "", "TLS mode (short)")

	flag.StringVar(&configFileShort, "c", "", "Path to the SMTP config file (short)")

	flag.StringVar(&fromEmailShort, "f", "", "Email address to send from (short)")
	flag.StringVar(&toEmailShort, "t", "", "Email addresses that will receive the email, comma-separated (short)")

	flag.StringVar(&subjectShort, "h", "", "Subject of the email (short)")
	flag.StringVar(&bodyShort, "b", "", "Body of the email (short)")

	flag.StringVar(&attachmentsFilesShort, "af", "", "File paths for attachments, comma-separated (short)")
	flag.StringVar(&bodyFileShort, "bf", "", "File path for email body (short)")
}

func main() {
	// Override the default flag.Usage
	flag.Usage = Usage
	flag.Parse()

	// Override long-form flags with short-form flags if set
	if smtpServerShort != "" {
		smtpServer = smtpServerShort
	}
	if smtpPortShort != 587 {
		smtpPort = smtpPortShort
	}
	if usernameShort != "" {
		username = usernameShort
	}
	if passwordShort != "" {
		password = passwordShort
	}
	if tlsModeShort != "" {
		tlsMode = tlsModeShort
	}
	if configFileShort != "" {
		configFile = configFileShort
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

	// Load config from file if provided
	if configFile != "" {
		config, err := loadConfig(configFile)
		if err != nil {
			fmt.Printf("Error loading config file: %v", err)
		}

		// Override flags with config file values if set
		if config.SMTPServer != "" {
			smtpServer = config.SMTPServer
		}
		if config.SMTPPort != 0 {
			smtpPort = config.SMTPPort
		}
		if config.SMTPUsername != "" {
			username = config.SMTPUsername
		}
		if config.SMTPPassword != "" {
			password = config.SMTPPassword
		}
		if config.TLSMode != "" {
			tlsMode = config.TLSMode
		}
		if config.FromEmail != "" {
			fromEmail = config.FromEmail
		}
	}

	// Check if required flags or config values are missing
	if smtpServer == "" || tlsMode == "" || fromEmail == "" || toEmail == "" || subject == "" {
		fmt.Fprintln(os.Stderr, "Error: Required flags or config values are missing.")
		Usage()
	}

	// Check if either direct input or file path is provided for body
	if body == "" && bodyFile == "" {
		fmt.Fprintln(os.Stderr, "Error: Subject and body are required, either directly or through a specified file.")
		Usage()
	}

	// Read body from files if provided
	if bodyFile != "" {
		content, err := os.ReadFile(bodyFile)
		if err != nil {
			fmt.Printf("\nError reading body file: %v\n", err)
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
		Usage()
	}

	sendEmail(smtpServer, smtpPort, username, password, fromEmail, toEmails, subject, body, bodyFile, attachmentPaths, tlsMode)
}
