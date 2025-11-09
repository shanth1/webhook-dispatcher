package github

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/shanth1/hookrelay/internal/common"
	"github.com/shanth1/hookrelay/internal/core/domain"
	"github.com/shanth1/hookrelay/internal/core/ports"
)

type Handler struct {
	templates *template.Template
	secret    string
}

var _ ports.WebhookHandler = (*Handler)(nil)

func NewHandler(secret string) (ports.WebhookHandler, error) {
	tmpls, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to parse github templates: %w", err)
	}
	return &Handler{templates: tmpls}, nil
}

func (h *Handler) Handle(ctx context.Context, req ports.WebhookRequest) (*domain.Notification, error) {
	if ok := h.verify(req); !ok {
		return nil, common.ErrInvalidSignature
	}

	payload, eventName, err := parsePayload(req)
	if err != nil {
		return nil, common.ErrInvalidSignature
	}

	templateName := common.GetTemplatePath(eventName)
	if h.templates.Lookup(templateName) == nil {
		templateName = common.GetTemplatePath("default")
	}

	var message bytes.Buffer
	if err := h.templates.ExecuteTemplate(&message, templateName, payload); err != nil {
		return nil, fmt.Errorf("error executing github template '%s': %w", templateName, err)
	}

	return &domain.Notification{Body: message.String()}, nil
}
