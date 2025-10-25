package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shanth1/gitrelay/internal/config"
)

var _ Sender = (*TelegramSender)(nil)

type TelegramSender struct {
	client *http.Client
	token  string
}

func NewTelegramSender(cfg config.TelegramSettings) *TelegramSender {
	return &TelegramSender{
		client: &http.Client{Timeout: 10 * time.Second},
		token:  cfg.Token,
	}
}

func (s *TelegramSender) Send(ctx context.Context, chatID string, text string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.token)

	payload := map[string]string{
		"chat_id": chatID,
		"text":    text,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal json payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: status %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}
