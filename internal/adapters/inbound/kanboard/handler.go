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
	URL       string                 `json:"url"`
	EventData map[string]interface{} `json:"event_data"`
}

type Handler struct {
	secret                  string
	baseURL                 string
	templates               *template.Template
	disableUnknownTemplates bool
}

var _ ports.WebhookHandler = (*Handler)(nil)

func NewHandler(secret, baseURL string, disableUnknownTemplates bool) (ports.WebhookHandler, error) {
	if secret == "" {
		return nil, fmt.Errorf("empty 'secret' value")
	}

	// TODO: added url validation
	if baseURL == "" {
		return nil, fmt.Errorf("empty 'base_url' value")
	}

	tmpls, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to parse kanboard templates: %w", err)
	}
	return &Handler{
		secret:                  secret,
		baseURL:                 baseURL,
		templates:               tmpls,
		disableUnknownTemplates: disableUnknownTemplates,
	}, nil
}

func (h *Handler) Handle(ctx context.Context, req ports.WebhookRequest) (*domain.Notification, error) {
	if ok := h.verify(req); !ok {
		return nil, common.ErrInvalidSignature
	}

	if req.Headers["content-type"] != "application/json" {
		return nil, fmt.Errorf("unsupported content type: %s", req.Headers["content-type"])
	}

	var payload KanboardPayload
	if err := json.Unmarshal(req.Payload, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kanboard json payload: %w", err)
	}

	if payload.EventName == "" {
		return nil, fmt.Errorf("kanboard event_name is missing from payload")
	}

	payload.EventData["eventName"] = payload.EventName

	// TODO: refactor
	if taskData, ok := payload.EventData["task"].(map[string]interface{}); ok {
		taskID, tid_ok := taskData["id"].(string)
		projectID, pid_ok := taskData["project_id"].(string)

		if tid_ok && pid_ok {
			taskURL := fmt.Sprintf("%s/?controller=TaskViewController&action=show&task_id=%s&project_id=%s", h.baseURL, taskID, projectID)
			taskData["url"] = taskURL
		}
	}

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
