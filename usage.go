package main

import (
	"fmt"
	"os"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  -s, --smtp-server         SMTP server for sending emails")
	fmt.Fprintln(os.Stderr, "  -p, --smtp-port           SMTP server port (Default: 587)")
	fmt.Fprintln(os.Stderr, "  -u, --smtp-username       Username for SMTP authentication")
	fmt.Fprintln(os.Stderr, "  -w, --smtp-password       Password for SMTP authentication")
	fmt.Fprintln(os.Stderr, "  -l, --tls-mode            TLS mode (none, tls-skip, tls)")
	fmt.Fprintln(os.Stderr, "  -f, --from-email          Email address to send from")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  -c, --config              Path to the SMTP json config file which replaces the above arguments")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  -t, --to-email            Email addresses that will receive the email, comma-separated")
	fmt.Fprintln(os.Stderr, "  -h, --subject             Subject of the email")
	fmt.Fprintln(os.Stderr, "  -b, --body                Body of the email")
	fmt.Fprintln(os.Stderr, "  -af, --attachments        File paths for attachments, comma-separated")
	fmt.Fprintln(os.Stderr, "  -bf, --body-file          File path for email body")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Ensure all required flags are provided.")
	os.Exit(1)
}
