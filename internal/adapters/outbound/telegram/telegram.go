package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/core/domain"
	"github.com/shanth1/hookrelay/internal/core/ports"
	"github.com/shanth1/hookrelay/internal/templates"
)

var _ ports.Notifier = (*Sender)(nil)

type Sender struct {
	client    *http.Client
	token     string
	templates *template.Template
}

func NewSender(cfg config.TelegramSettings, registry *templates.Registry) (ports.Notifier, error) {
	funcMap := template.FuncMap{
		"escape": escapeMarkdown,
	}

	allTemplates, err := registry.LoadAll(funcMap)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates for telegram sender: %w", err)
	}

	return &Sender{
		client:    &http.Client{Timeout: 10 * time.Second},
		token:     cfg.Token,
		templates: allTemplates,
	}, nil
}

func (s *Sender) Send(ctx context.Context, chatID string, notification domain.Notification) error {
	var textToSend string
	var parseMode = ""

	if notification.TemplateName != "" && notification.TemplateData != nil {
		var buf bytes.Buffer
		err := s.templates.ExecuteTemplate(&buf, notification.TemplateName, notification.TemplateData)
		if err != nil {
			return fmt.Errorf("failed to execute telegram template '%s': %w", notification.TemplateName, err)
		}
		textToSend = buf.String()
		parseMode = "MarkdownV2"
	} else {
		textToSend = notification.Body
	}

	err := s.trySend(ctx, chatID, textToSend, parseMode)
	if err != nil && strings.Contains(err.Error(), "can't parse entities") {
		log.FromContext(ctx).Warn().Str("template", notification.TemplateName).Msg("MarkdownV2 parsing failed, falling back to plain text.")
		return s.trySend(ctx, chatID, textToSend, "") // Повторная отправка без форматирования
	}

	return err
}

func (s *Sender) trySend(ctx context.Context, chatID, text, parseMode string) error {
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

func escapeMarkdown(text interface{}) string {
	if text == nil {
		return ""
	}

	s := fmt.Sprint(text)

	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(s)
}
