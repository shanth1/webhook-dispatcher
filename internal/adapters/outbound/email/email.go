package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/core/domain"
	"github.com/shanth1/hookrelay/internal/core/ports"
)

var _ ports.Notifier = (*Sender)(nil)

type Sender struct {
	cfg  config.EmailSettings
	auth smtp.Auth
}

func NewSender(cfg config.EmailSettings) *Sender {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return &Sender{
		cfg:  cfg,
		auth: auth,
	}
}

func (s *Sender) Send(ctx context.Context, recipientEmail string, notification domain.Notification) error {
	subject := "Webhook Notification"
	msg := []byte("From: " + s.cfg.From + "\r\n" +
		"To: " + recipientEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		notification.Body)

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
