<p align="center">
 <img src="https://github.com/KeepSec-Technologies/Mail2Go/assets/108779415/f28159ee-5300-45bd-9e53-76ca2bf2ca88"
</p>

# Mail2Go - Lightweight CLI SMTP client

Mail2Go is a very lightweight command-line SMTP client written in Go, designed to send emails from the command-line easily.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Building from Source](#building-from-source)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Send Emails with Ease**: Quickly send emails with subject, body, and multiple recipients.
- **Attachments Support**: Attach multiple files of various types to your emails.
- **HTML and Plain Text**: Supports both HTML and plain text formats for email bodies.
- **Command Line Interface**: Easy-to-use CLI arguments for configuring and sending emails.
- **Flexible Configuration**: SMTP server, port, username, and password can be configured through CLI arguments.

## Requirements

- Go 1.20 or higher recommended (for build).
- Access to an SMTP server for sending emails.

## Installation

1. Download the binary with wget:

    ```shell
    wget https://github.com/KeepSec-Technologies/Mail2Go/releases/download/1.0/mail2go_linux_amd64_1.0.tar.gz
    ```

2. Unpack it with tar

    ```shell
    tar -xf mail2go_linux_amd64_1.0.tar.gz
    ```

3. Move it to your /usr/local/bin/ (Optional):

    ```shell
    sudo mv mail2go /usr/local/bin/mail2go
    ```

## Building from Source

1. Ensure you have Go installed on your system. You can download Go from [here](https://golang.org/dl/).
2. Clone the repository:

    ```shell
    git clone https://github.com/KeepSec-Technologies/Mail2Go
    ```

3. Navigate to the cloned directory:

    ```shell
    cd Mail2Go
    ```

4. Build the tool:

    ```shell
    CGO_ENABLED=0 go build -a -installsuffix cgo -o mail2go .
    ```

## Usage

Run the Mail2Go tool with the required flags:

```shell
./mail2go --smtp-server [SMTP_SERVER] --smtp-port [SMTP_PORT] --smtp-username [USERNAME] --smtp-password [PASSWORD] --from-email [FROM_EMAIL] --to-email [TO_EMAIL_1],[TO_EMAIL_2] --subject "Your Subject" --body "Your email body." --attachments "path/to/attachment1,path/to/attachment2"
```

Flags:

```text
 -s, --smtp-server         SMTP server for sending emails
 -p, --smtp-port           SMTP server port (Default: 587)
 -u, --smtp-username       Username for SMTP authentication
 -w, --smtp-password       Password for SMTP authentication
 -f, --from-email          Email address to send from
 -t, --to-email            Email addresses that will receive the email, comma-separated
 -h, --subject             Subject of the email
 -b, --body                Body of the email
 -af, --attachments        File paths for attachments, comma-separated
 -bf, --body-file          File path for email body
```

Basic example:

```shell
./mail2go -s mail.example.com -u user@example.com -w password123 -f mail2go@example.com -t admin@example.com -h 'Test Mail2Go Subject' -b 'This is a body!' 
```

Example with two recipients, the body from an HTML file and two attached files (can be more):

```shell
./mail2go -s mail.example.com -u user@example.com -w password123 -f mail2go@example.com -t admin@example.com,other@example.com -h 'Test Mail2Go Subject' -bf demo/body.html -af README.md,demo/mail2go-smaller.png
```

## Contributing

Contributions to Mail2Go are welcome. Please fork the repository and submit a pull request with your changes or improvements.

## License

This project is licensed under MIT - see the LICENSE file for details.
