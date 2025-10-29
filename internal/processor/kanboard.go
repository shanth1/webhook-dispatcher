package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/shanth1/hookrelay/internal/utils"
)

type KanboardPayload struct {
	EventName string                 `json:"event_name"`
	EventData map[string]interface{} `json:"event_data"`
}

type KanboardProcessor struct {
	templates *template.Template
}

func NewKanboardProcessor(templates *template.Template) *KanboardProcessor {
	return &KanboardProcessor{templates: templates}
}

func (p *KanboardProcessor) Process(r *http.Request) (string, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return "", fmt.Errorf("unsupported content type: %s", r.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read request body: %w", err)
	}

	var payload KanboardPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("failed to unmarshal kanboard json payload: %w", err)
	}

	if payload.EventName == "" {
		return "", fmt.Errorf("kanboard event_name is missing from payload")
	}

	payload.EventData["eventName"] = payload.EventName

	templateName := utils.GetTemplatePath("kanboard", payload.EventName)
	if p.templates.Lookup(templateName) == nil {
		templateName = utils.GetTemplatePath("kanboard", "default")
		if p.templates.Lookup(templateName) == nil {
			return "", fmt.Errorf("default kanboard template '%s' not found", templateName)
		}
	}

	var message bytes.Buffer
	err = p.templates.ExecuteTemplate(&message, templateName, payload)
	if err != nil {
		return "", fmt.Errorf("error executing kanboard template '%s': %w", templateName, err)
	}

	return message.String(), nil
}
