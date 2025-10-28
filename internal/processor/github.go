package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"

	"github.com/shanth1/gitrelay/internal/utils"
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
	r.Body = io.NopCloser(bytes.NewBuffer(body))

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

	var message bytes.Buffer
	tmpl := p.templates.Lookup(utils.GetTemplatePath("github", eventName))
	if tmpl == nil {
		tmpl = p.templates.Lookup(utils.GetTemplatePath("github", "default"))
	}
	if err := tmpl.Execute(&message, payload); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return message.String(), nil
}
