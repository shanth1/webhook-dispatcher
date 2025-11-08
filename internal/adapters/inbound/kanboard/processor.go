package kanboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"path"

	"github.com/shanth1/hookrelay/internal/core/domain"
	"github.com/shanth1/hookrelay/internal/service/webhook"
)

type KanboardPayload struct {
	EventName string                 `json:"event_name"`
	EventData map[string]interface{} `json:"event_data"`
}

type Processor struct {
	templates *template.Template
}

var _ webhook.Processor = (*Processor)(nil)

func NewProcessor() (webhook.Processor, error) {
	tmpls, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to parse kanboard templates: %w", err)
	}
	return &Processor{templates: tmpls}, nil
}

func (p *Processor) Process(body []byte, headers map[string]string) (*domain.Notification, error) {
	if headers["Content-Type"] != "application/json" {
		return nil, fmt.Errorf("unsupported content type: %s", headers["Content-Type"])
	}

	var payload KanboardPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kanboard json payload: %w", err)
	}

	if payload.EventName == "" {
		return nil, fmt.Errorf("kanboard event_name is missing from payload")
	}

	payload.EventData["eventName"] = payload.EventName

	templateName := getTemplatePath(payload.EventName)
	if p.templates.Lookup(templateName) == nil {
		templateName = getTemplatePath("default")
	}

	var message bytes.Buffer
	if err := p.templates.ExecuteTemplate(&message, templateName, payload.EventData); err != nil {
		return nil, fmt.Errorf("error executing kanboard template '%s': %w", templateName, err)
	}

	return &domain.Notification{Body: message.String()}, nil
}

func getTemplatePath(name string) string {
	return path.Join("templates", fmt.Sprintf("%s.tmpl", name))
}
