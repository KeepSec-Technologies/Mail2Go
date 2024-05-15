# Use an official Debian image as a parent image
FROM debian:12

# Set noninteractive installation to avoid prompts
ENV DEBIAN_FRONTEND=noninteractive

# Install Postfix
RUN apt-get update && apt-get install -y postfix

# Copy your TLS certificates into the container
COPY smtp.crt /etc/ssl/certs/smtp.crt
COPY smtp.key /etc/ssl/private/smtp.key

# Setup Postfix configuration
RUN postconf -e "myhostname = mydomain.com" \
    && postconf -e "smtpd_tls_cert_file = /etc/ssl/certs/smtp.crt" \
    && postconf -e "smtpd_tls_key_file = /etc/ssl/private/smtp.key" \
    && postconf -e "smtpd_tls_security_level = may" \
    && postconf -e "smtpd_tls_auth_only = no" \
    && postconf -e "smtpd_relay_restrictions = permit_mynetworks permit_sasl_authenticated defer_unauth_destination" \
    && postconf -e "mydestination = $myhostname, localhost.$mydomain, localhost" \
    && postconf -e "mynetworks = 127.0.0.0/8 [::ffff:127.0.0.0]/104 [::1]/128 0.0.0.0/0" \
    && postconf -e "inet_interfaces = all" \
    && postconf -e "inet_protocols = ipv4" \
    && postconf -e "smtpd_recipient_restrictions = permit_mynetworks permit_sasl_authenticated reject_unauth_destination" \
    && postconf -e "smtp_bind_address = 0.0.0.0" \
    && postconf -e "smtpd_sasl_auth_enable = no"

# Expose ports for SMTP
EXPOSE 25

# Start Postfix
CMD ["postfix", "start-fg"]

