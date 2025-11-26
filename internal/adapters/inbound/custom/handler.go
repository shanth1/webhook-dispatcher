package custom

import (
	"context"
	"fmt"

	"github.com/shanth1/hookrelay/internal/common"
	"github.com/shanth1/hookrelay/internal/core/domain"
	"github.com/shanth1/hookrelay/internal/core/ports"
)

type Handler struct {
	secret string
}

var _ ports.WebhookHandler = (*Handler)(nil)

func NewHandler(secret string) ports.WebhookHandler {
	return &Handler{
		secret: secret,
	}
}

func (h *Handler) Handle(ctx context.Context, req ports.WebhookRequest) (*domain.Notification, error) {
	if ok := h.verify(req); !ok {
		return nil, common.ErrInvalidSignature
	}

	if len(req.Payload) == 0 {
		return nil, fmt.Errorf("request payload is empty")
	}

	return &domain.Notification{Body: string(req.Payload)}, nil
}

func (h *Handler) verify(req ports.WebhookRequest) bool {
	providerSecret := req.GetHeader("X-Auth-Token")

	return providerSecret != "" && providerSecret == h.secret
}
