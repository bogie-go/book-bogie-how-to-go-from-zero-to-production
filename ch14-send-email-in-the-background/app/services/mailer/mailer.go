// Package mailer sends email over SMTP.
package mailer

import (
	"context"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"

	"depot/config"
)

// Message is one email.
type Message struct {
	To      string
	Subject string
	Body    string
}

// SMTP sends messages through the server in the app's SMTP configuration.
type SMTP struct {
	cfg config.SmtpConfig
}

// New builds a sender from the app's SMTP configuration.
func New(cfg config.SmtpConfig) *SMTP {
	return &SMTP{cfg: cfg}
}

// Send delivers msg, or returns why it couldn't.
func (s *SMTP) Send(ctx context.Context, msg Message) error {
	from, err := mail.ParseAddress(s.cfg.From)
	if err != nil {
		return fmt.Errorf("mailer: the sender %q isn't an email address: %w", s.cfg.From, err)
	}
	addr := net.JoinHostPort(s.cfg.Address, strconv.Itoa(s.cfg.Port))
	var auth smtp.Auth
	if s.cfg.Username != "" {
		auth = smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Address)
	}
	body := strings.Join([]string{
		"From: " + s.cfg.From,
		"To: " + msg.To,
		"Subject: " + msg.Subject,
		"Content-Type: text/plain; charset=utf-8",
		"",
		msg.Body,
	}, "\r\n")
	if err := smtp.SendMail(addr, auth, from.Address, []string{msg.To}, []byte(body)); err != nil {
		return fmt.Errorf("mailer: send to %s: %w", msg.To, err)
	}
	return nil
}
