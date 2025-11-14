package github

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/shanth1/hookrelay/internal/core/ports"
)

func parsePayload(req ports.WebhookRequest) (payload map[string]interface{}, eventName string, err error) {
	payloadJSON := ""
	contentType := req.Headers["content-type"]

	if contentType == "application/x-www-form-urlencoded" {
		form, err := url.ParseQuery(string(req.Payload))
		if err != nil {
			return nil, "", fmt.Errorf("error parsing form payload: %w", err)
		}
		payloadJSON = form.Get("payload")
	} else if contentType == "application/json" {
		payloadJSON = string(req.Payload)
	}

	if payloadJSON == "" {
		return nil, "", fmt.Errorf("payload is empty or content type is unsupported")
	}

	eventName = req.Headers["x-github-event"]
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		return nil, "", fmt.Errorf("error parsing JSON payload: %w", err)
	}
	payload["eventName"] = eventName

	return
}
