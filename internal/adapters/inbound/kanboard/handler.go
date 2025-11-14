package kanboard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/shanth1/hookrelay/internal/common"
	"github.com/shanth1/hookrelay/internal/core/domain"
	"github.com/shanth1/hookrelay/internal/core/ports"
)

type KanboardPayload struct {
	EventName string                 `json:"event_name"`
	EventData map[string]interface{} `json:"event_data"`
}

type Handler struct {
	secret                  string
	templates               *template.Template
	disableUnknownTemplates bool
}

var _ ports.WebhookHandler = (*Handler)(nil)

func NewHandler(secret string, disableUnknownTemplates bool) (ports.WebhookHandler, error) {
	tmpls, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to parse kanboard templates: %w", err)
	}
	return &Handler{
		secret:                  secret,
		templates:               tmpls,
		disableUnknownTemplates: disableUnknownTemplates,
	}, nil
}

func (h *Handler) Handle(ctx context.Context, req ports.WebhookRequest) (*domain.Notification, error) {
	if ok := h.verify(req); !ok {
		return nil, common.ErrInvalidSignature
	}

	if req.Headers["Content-Type"] != "application/json" {
		return nil, fmt.Errorf("unsupported content type: %s", req.Headers["Content-Type"])
	}

	var payload KanboardPayload
	if err := json.Unmarshal(req.Payload, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kanboard json payload: %w", err)
	}

	if payload.EventName == "" {
		return nil, fmt.Errorf("kanboard event_name is missing from payload")
	}

	payload.EventData["eventName"] = payload.EventName

	templateName := common.GetTemplatePath(payload.EventName)
	if h.templates.Lookup(templateName) == nil {
		templateExists := h.templates.Lookup(templateName) != nil
		if !templateExists {
			if h.disableUnknownTemplates {
				return nil, nil
			}
			templateName = common.GetTemplatePath("default")
		}
	}

	var message bytes.Buffer
	if err := h.templates.ExecuteTemplate(&message, templateName, payload.EventData); err != nil {
		return nil, fmt.Errorf("error executing kanboard template '%s': %w", templateName, err)
	}

	return &domain.Notification{Body: message.String()}, nil
}
