<p align="center">
 <img src="https://github.com/KeepSec-Technologies/Mail2Go/assets/108779415/afb750ff-0320-46d5-8e03-e26dfbde1c49"
</p>

# Mail2Go - Lightweight CLI SMTP client

![License](https://img.shields.io/github/license/KeepSec-Technologies/Mail2Go)
![GitHub issues](https://img.shields.io/github/issues-raw/KeepSec-Technologies/Mail2Go)
![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/KeepSec-Technologies/Mail2Go/main)

Mail2Go is a very lightweight command-line SMTP client written in Go, designed to send emails from the command-line easily.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Building from Source](#building-from-source)
- [Usage](#usage)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Send Emails with Ease**: Quickly send emails with subject, body, and multiple recipients.
- **Attachments Support**: Attach multiple files of various types to your emails.
- **HTML and Plain Text**: Supports both HTML and plain text formats for email bodies.
- **Command Line Interface**: Easy-to-use CLI arguments for configuring and sending emails.
- **Flexible Configuration**: SMTP server, TLS, port, username, and password can be configured through CLI arguments or a JSON configuration file.

## Requirements

- Go 1.20 or higher recommended (for build).
- Access to an SMTP server for sending emails.

## Installation

1. Download the Linux amd64 binary with wget (more versions on [release](https://github.com/KeepSec-Technologies/Mail2Go/releases/tag/1.1.4) tab):

    ```shell
    wget https://github.com/KeepSec-Technologies/Mail2Go/releases/download/1.1.4/mail2go_linux_amd64_1.1.4.tar.gz
    ```

2. Unpack it with tar

    ```shell
    tar -xf mail2go_linux_amd64_1.1.4.tar.gz
    ```

3. Move it to your /usr/local/bin/ (Optional):

    ```shell
    sudo mv mail2go /usr/local/bin/mail2go
    ```

## Building from Source

1. Ensure you have Go installed on your system. You can download Go from [here](https://go.dev/dl/).
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

Run the Mail2Go tool with the required arguments:

```text
-s, --smtp-server         SMTP server for sending emails
-p, --smtp-port           SMTP server port (Default: 587)
-u, --smtp-username       Username for SMTP authentication
-w, --smtp-password       Password for SMTP authentication
-l, --tls-mode            TLS mode (none, tls-skip, tls)
-f, --from-email          Email address to send from

-c, --config              Path to the SMTP json config file which replaces the above arguments

-t, --to-email            Email addresses that will receive the email, can be multiples (comma-separated)
-h, --subject             Subject of the email
-b, --body                Body of the email
-af, --attachments        File paths for attachments, can be multiples (comma-separated)
-bf, --body-file          File path for email body
```

## Examples

```shell
# Basic example:
mail2go -s mail.example.com -p 587 -u user@example.com -w password123 -l tls -f mail2go@example.com -t admin@example.com -h 'Test Mail2Go Subject' -b 'This is a body!' 

# Example with two recipients, the body from an HTML file and two attached files (can be more):
mail2go -s mail.example.com -p 587  -u user@example.com -w password123 -l tls -f mail2go@example.com -t admin@example.com,other@example.com -h 'Test Mail2Go Subject' -bf demo/body.html -af README.md,demo/mail2go-smaller.png

# Example without authentication and no TLS:
mail2go -s mail.example.com -p 25 -l none -f mail2go@example.com -t admin@example.com -h 'Test Mail2Go Subject' -b 'This is a body!' 

# Example of a local SMTP server with a TLS certificate that can't be verified:
mail2go -s localhost -p 587 -u user@example.com -w password123 -l tls-skip -f mail2go@example.com -t admin@example.com -h 'Test Mail2Go Subject' -b 'This is a body!' 
```

Example of json configuration file to pass to Mail2Go:

```json
{
    "smtp_server": "mail.example.com",
    "smtp_port": 587,
    "smtp_username": "user@example.com",
    "smtp_password": "password_example",
    "tls_mode": "tls",
    "from_email": "mail2go@example.com"
}
```

Example of the command that goes with it:

```shell
mail2go -c demo/config.json -t admin@example.com -h 'Test Mail2Go Subject' -b 'This is a body!' 
```

## Contributing

Contributions to Mail2Go are welcome. Please fork the repository and submit a pull request with your changes or improvements.

## License

This project is licensed under MIT - see the LICENSE file for details.
