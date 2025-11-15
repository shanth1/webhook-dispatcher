package kanboard

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/shanth1/hookrelay/internal/common"
	"github.com/shanth1/hookrelay/internal/core/domain"
	"github.com/shanth1/hookrelay/internal/core/ports"
	"github.com/shanth1/hookrelay/internal/templates"
)

type KanboardPayload struct {
	EventName string                 `json:"event_name"`
	URL       string                 `json:"url"`
	EventData map[string]interface{} `json:"event_data"`
}

type Handler struct {
	secret                  string
	baseURL                 string
	templateFS              fs.FS
	disableUnknownTemplates bool
}

var _ ports.WebhookHandler = (*Handler)(nil)

func NewHandler(secret, baseURL string, disableUnknownTemplates bool, registry *templates.Registry) (ports.WebhookHandler, error) {
	if secret == "" {
		return nil, fmt.Errorf("empty 'secret' value")
	}

	// TODO: added url validation
	if baseURL == "" {
		return nil, fmt.Errorf("empty 'base_url' value")
	}

	err := registry.RegisterSource("kanboard", templates.Source{
		FS:       templateFiles,
		Patterns: []string{"templates/*.tmpl"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register kanboard template source: %w", err)
	}

	return &Handler{
		secret:                  secret,
		baseURL:                 baseURL,
		templateFS:              templateFiles,
		disableUnknownTemplates: disableUnknownTemplates,
	}, nil
}

func (h *Handler) Handle(ctx context.Context, req ports.WebhookRequest) (*domain.Notification, error) {
	if ok := h.verify(req); !ok {
		return nil, common.ErrInvalidSignature
	}

	if req.GetHeader("Content-Type") != "application/json" {
		return nil, fmt.Errorf("unsupported content type: %s", req.GetHeader("Content-Type"))
	}

	var payload KanboardPayload
	if err := json.Unmarshal(req.Payload, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kanboard json payload: %w", err)
	}

	if payload.EventName == "" {
		return nil, fmt.Errorf("kanboard event_name is missing from payload")
	}

	// TODO: refactor
	if taskData, ok := payload.EventData["task"].(map[string]interface{}); ok {
		taskIDFloat, tid_ok := taskData["id"].(float64)
		projectIDFloat, pid_ok := taskData["project_id"].(float64)

		if tid_ok && pid_ok {
			taskURL := fmt.Sprintf(
				"%s/?controller=TaskViewController&action=show&task_id=%d&project_id=%d",
				h.baseURL,
				int64(taskIDFloat),
				int64(projectIDFloat),
			)
			taskData["url"] = taskURL
		}
	}

	templateName := common.GetTemplatePath(payload.EventName)
	if _, err := fs.Stat(h.templateFS, "templates/"+templateName); err != nil {
		if h.disableUnknownTemplates {
			return nil, nil
		}
		templateName = common.GetTemplatePath("default")
	}

	return &domain.Notification{
		TemplateName: "kanboard/" + templateName,
		TemplateData: payload,
	}, nil
}
