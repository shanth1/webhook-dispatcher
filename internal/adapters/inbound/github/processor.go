package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"path"

	"github.com/shanth1/hookrelay/internal/core/domain"
	"github.com/shanth1/hookrelay/internal/service/webhook"
)

type Processor struct {
	templates *template.Template
}

var _ webhook.Processor = (*Processor)(nil)

func NewProcessor() (webhook.Processor, error) {
	tmpls, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to parse github templates: %w", err)
	}
	return &Processor{templates: tmpls}, nil
}

func (p *Processor) Process(body []byte, headers map[string]string) (*domain.Notification, error) {
	payloadJSON := ""
	contentType := headers["Content-Type"]

	if contentType == "application/x-www-form-urlencoded" {
		form, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, fmt.Errorf("error parsing form payload: %w", err)
		}
		payloadJSON = form.Get("payload")
	} else if contentType == "application/json" {
		payloadJSON = string(body)
	}

	if payloadJSON == "" {
		return nil, fmt.Errorf("payload is empty or content type is unsupported")
	}

	eventName := headers["X-GitHub-Event"]
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		return nil, fmt.Errorf("error parsing JSON payload: %w", err)
	}
	payload["eventName"] = eventName

	templateName := getTemplatePath(eventName)
	if p.templates.Lookup(templateName) == nil {
		templateName = getTemplatePath("default")
	}

	var message bytes.Buffer
	if err := p.templates.ExecuteTemplate(&message, templateName, payload); err != nil {
		return nil, fmt.Errorf("error executing github template '%s': %w", templateName, err)
	}

	return &domain.Notification{Body: message.String()}, nil
}

func getTemplatePath(name string) string {
	return path.Join("templates", fmt.Sprintf("%s.tmpl", name))
}
