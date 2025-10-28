package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"

	"github.com/shanth1/hookrelay/internal/utils"
)

type GithubProcessor struct {
	templates *template.Template
}

func NewGithubProcessor(templates *template.Template) *GithubProcessor {
	return &GithubProcessor{templates: templates}
}

func (p *GithubProcessor) Process(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read request body: %w", err)
	}

	payloadJSON := ""
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/x-www-form-urlencoded" {
		form, err := url.ParseQuery(string(body))
		if err != nil {
			return "", fmt.Errorf("error parsing form payload: %w", err)
		}
		payloadJSON = form.Get("payload")
	} else if contentType == "application/json" {
		payloadJSON = string(body)
	}

	if payloadJSON == "" {
		return "", fmt.Errorf("payload is empty or content type is unsupported")
	}

	eventName := r.Header.Get("X-GitHub-Event")
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		return "", fmt.Errorf("error parsing JSON payload: %w", err)
	}
	payload["eventName"] = eventName

	templateName := utils.GetTemplatePath("github", eventName)

	if p.templates.Lookup(templateName) == nil {
		templateName = utils.GetTemplatePath("github", "default")
		if p.templates.Lookup(templateName) == nil {
			return "", fmt.Errorf("default template '%s' not found", templateName)
		}
	}

	var message bytes.Buffer
	err = p.templates.ExecuteTemplate(&message, templateName, payload)
	if err != nil {
		return "", fmt.Errorf("error executing template '%s': %w", templateName, err)
	}

	return message.String(), nil
}
