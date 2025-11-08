package ports

import (
	"context"

	"github.com/shanth1/hookrelay/internal/config"
)

type WebhookRequest struct {
	WebhookType config.WebhookType
	Payload     []byte
	Headers     map[string]string
}

type WebhookService interface {
	ProcessWebhook(ctx context.Context, req WebhookRequest, recipients []config.Recipient) error
}
