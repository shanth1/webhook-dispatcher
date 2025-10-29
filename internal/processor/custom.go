package processor

import (
	"fmt"
	"io"
	"net/http"
)

type CustomProcessor struct{}

func NewCustomProcessor() *CustomProcessor {
	return &CustomProcessor{}
}

func (p *CustomProcessor) Process(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read request body: %w", err)
	}
	if len(body) == 0 {
		return "", fmt.Errorf("request body is empty")
	}
	return string(body), nil
}
