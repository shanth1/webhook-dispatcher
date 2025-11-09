package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"

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

	payloadJSON := ""
	contentType := req.Headers["Content-Type"]

	if contentType == "application/x-www-form-urlencoded" {
		form, err := url.ParseQuery(string(req.Payload))
		if err != nil {
			return nil, fmt.Errorf("error parsing form payload: %w", err)
		}
		payloadJSON = form.Get("payload")
	} else if contentType == "application/json" {
		payloadJSON = string(req.Payload)
	}

	if payloadJSON == "" {
		return nil, fmt.Errorf("payload is empty or content type is unsupported")
	}

	eventName := req.Headers["X-GitHub-Event"]
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		return nil, fmt.Errorf("error parsing JSON payload: %w", err)
	}
	payload["eventName"] = eventName

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
