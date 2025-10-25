package service

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/shanth1/gitrelay/internal/config"
)

var _ Sender = (*EmailSender)(nil)

type EmailSender struct {
	cfg  config.EmailSettings
	auth smtp.Auth
}

func NewEmailSender(cfg config.EmailSettings) *EmailSender {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return &EmailSender{
		cfg:  cfg,
		auth: auth,
	}
}

func (s *EmailSender) Send(ctx context.Context, recipientEmail string, message string) error {
	subject := "GitHub Notification"
	msg := []byte("From: " + s.cfg.From + "\r\n" +
		"To: " + recipientEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message)

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	errChan := make(chan error, 1)
	go func() {
		errChan <- smtp.SendMail(addr, s.auth, s.cfg.From, []string{recipientEmail}, msg)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("failed to send email via SMTP: %w", err)
		}
		return nil
	}
}
