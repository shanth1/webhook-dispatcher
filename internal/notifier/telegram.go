package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/config"
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
	err := s.trySend(ctx, chatID, text, "MarkdownV2")
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "can't parse entities") {
		log.FromContext(ctx).Warn().Msg("MarkdownV2 parsing failed, falling back to plain text.")

		return s.trySend(ctx, chatID, text, "")
	}

	return err
}

func (s *TelegramSender) trySend(ctx context.Context, chatID, text, parseMode string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.token)

	payload := map[string]string{
		"chat_id": chatID,
		"text":    text,
	}
	if parseMode != "" {
		payload["parse_mode"] = parseMode
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

// TODO: remove or implement
func escapeMarkdownV2(text string) string {
	charsToEscape := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

	for _, char := range charsToEscape {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}
