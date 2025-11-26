package ports

import (
	"context"
	"strings"

	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/core/domain"
)

type WebhookRequest struct {
	Payload []byte
	Headers map[string]string
	Params  map[string]string
}

func (r *WebhookRequest) GetHeader(key string) string {
	return r.Headers[strings.ToLower(key)]
}

type WebhookHandler interface {
	Handle(ctx context.Context, req WebhookRequest) (*domain.Notification, error)
}

type Service interface {
	ProcessWebhook(ctx context.Context, webhookName config.WebhookName, req WebhookRequest, recipients []config.Recipient) error
}
