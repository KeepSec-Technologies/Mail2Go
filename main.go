package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var version = "1.1.8"

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
	replyTo string

	subject string
	body    string

	attachmentsFiles string
	bodyFile         string

	showVersion bool

	// Short-form flags
	smtpServerShort string
	smtpPortShort   int
	usernameShort   string
	passwordShort   string

	noAuth      bool
	noAuthShort bool

	tlsModeShort string

	configFileShort string

	fromEmailShort string
	toEmailShort   string
	replyToShort string

	subjectShort string
	bodyShort    string

	attachmentsFilesShort string
	bodyFileShort         string

	showVersionShort bool
)

func init() {
	// Long-form flags
	flag.StringVar(&smtpServer, "smtp-server", "", "SMTP server for sending emails")
	flag.IntVar(&smtpPort, "smtp-port", 587, "SMTP server port")
	flag.StringVar(&username, "smtp-username", "", "Username for SMTP authentication")
	flag.StringVar(&password, "smtp-password", "", "Password for SMTP authentication")
	flag.BoolVar(&noAuth, "no-auth", false, "Use unauthenticated SMTP")

	flag.StringVar(&tlsMode, "tls-mode", "tls", "TLS mode (none, tls-skip, tls)")

	flag.StringVar(&configFile, "config", "", "Path to the SMTP config file")

	flag.StringVar(&fromEmail, "from-email", "", "Email address to send from")
	flag.StringVar(&toEmail, "to-email", "", "Email addresses that will receive the email, comma-separated")
	flag.StringVar(&replyTo, "reply-to", "", "Email address to reply to")

	flag.StringVar(&subject, "subject", "", "Subject of the email")
	flag.StringVar(&body, "body", "", "Body of the email")

	flag.StringVar(&attachmentsFiles, "attachments", "", "File paths for attachments, comma-separated")
	flag.StringVar(&bodyFile, "body-file", "", "File path for email body")

	flag.BoolVar(&showVersion, "version", false, "Display application version")

	// Short-form flags
	flag.StringVar(&smtpServerShort, "s", "", "SMTP server for sending emails (short)")
	flag.IntVar(&smtpPortShort, "p", 587, "SMTP server port (short)")
	flag.StringVar(&usernameShort, "u", "", "Username for SMTP authentication (short)")
	flag.StringVar(&passwordShort, "w", "", "Password for SMTP authentication (short)")
	flag.BoolVar(&noAuthShort, "na", false, "Use unauthenticated SMTP (short)")

	flag.StringVar(&tlsModeShort, "l", "tls", "TLS mode (short)")

	flag.StringVar(&configFileShort, "c", "", "Path to the SMTP config file (short)")

	flag.StringVar(&fromEmailShort, "f", "", "Email address to send from (short)")
	flag.StringVar(&toEmailShort, "t", "", "Email addresses that will receive the email, comma-separated (short)")
	flag.StringVar(&replyToShort, "r", "", "Email address to reply to (short)")

	flag.StringVar(&subjectShort, "h", "", "Subject of the email (short)")
	flag.StringVar(&bodyShort, "b", "", "Body of the email (short)")

	flag.StringVar(&attachmentsFilesShort, "af", "", "File paths for attachments, comma-separated (short)")
	flag.StringVar(&bodyFileShort, "bf", "", "File path for email body (short)")

	flag.BoolVar(&showVersionShort, "v", false, "Display application version")
}

func main() {
	// Override the default flag.Usage
	flag.Usage = Usage
	flag.Parse()

	showVersion = showVersion || showVersionShort
	if showVersion {
		fmt.Printf("Mail2Go Version: %s\n", version)
		os.Exit(0)
	}

	var config Config = Config{}

	// Load config file
	configFile = priorityString([]string{configFile, configFileShort})
	if configFile == "" {
		// config file not provided -- look for default file
		if path, err := os.UserConfigDir(); err == nil {
			path = filepath.Join(path, "mail2go", "config.json")
			if _, err := os.Stat(path); err == nil {
				configFile = path
			}
		}
	}
	if configFile != "" {
		c, err := loadConfig(configFile)
		if err != nil {
			fmt.Printf("Error loading config file: %v", err)
		}
		config = c
	}

	// Clearly define our config priorities, lowest to highest: config files, long flags, short flags
	smtpServer = priorityString([]string{config.SMTPServer, smtpServer, smtpServerShort})
	smtpPort = priorityInt(587, []int{config.SMTPPort, smtpPort, smtpPortShort})
	username = priorityString([]string{config.SMTPUsername, username, usernameShort})
	password = priorityString([]string{config.SMTPPassword, password, passwordShort})
	noAuth = config.NoAuth || noAuth || noAuthShort
	tlsMode = priorityString([]string{config.TLSMode, tlsMode, tlsModeShort})
	fromEmail = priorityString([]string{config.FromEmail, fromEmail, fromEmailShort})

	toEmail = priorityString([]string{toEmail, toEmailShort})
	replyTo = priorityString([]string{replyTo, replyToShort})
	subject = priorityString([]string{subject, subjectShort})
	body = priorityString([]string{body, bodyShort})
	attachmentsFiles = priorityString([]string{attachmentsFiles, attachmentsFilesShort})
	bodyFile = priorityString([]string{bodyFile, bodyFileShort})

	// Check if required flags or config values are missing
	if smtpServer == "" || fromEmail == "" || toEmail == "" || subject == "" {
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
		body = priorityString([]string{string(content), body}) //preserve the "flags override files" semantic
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

	sendEmail(smtpServer, smtpPort, username, password, fromEmail, toEmails, replyTo, subject, body, bodyFile, attachmentPaths, tlsMode, noAuth)
}

func priorityString(strings []string) string {
	var result = ""
	for _, val := range strings {
		if val != "" {
			result = val
		}
	}
	return result
}

func priorityInt(emptyval int, ints []int) int {
	var result = emptyval
	for _, val := range ints {
		if val != emptyval {
			result = val
		}
	}
	return result
}
