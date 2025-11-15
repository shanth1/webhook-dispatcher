package github

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/shanth1/hookrelay/internal/common"
	"github.com/shanth1/hookrelay/internal/core/domain"
	"github.com/shanth1/hookrelay/internal/core/ports"
	"github.com/shanth1/hookrelay/internal/templates"
)

type Handler struct {
	secret                  string
	templateFS              fs.FS
	disableUnknownTemplates bool
}

var _ ports.WebhookHandler = (*Handler)(nil)

func NewHandler(secret string, disableUnknownTemplates bool, registry *templates.Registry) (ports.WebhookHandler, error) {
	err := registry.RegisterSource("github", templates.Source{
		FS:       templateFiles,
		Patterns: []string{"templates/*.tmpl"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register github template source: %w", err)
	}

	return &Handler{
		secret:                  secret,
		templateFS:              templateFiles,
		disableUnknownTemplates: disableUnknownTemplates,
	}, nil
}

func (h *Handler) Handle(ctx context.Context, req ports.WebhookRequest) (*domain.Notification, error) {
	if ok := h.verify(req); !ok {
		return nil, common.ErrInvalidSignature
	}

	payload, eventName, err := parsePayload(req)
	if err != nil {
		return nil, fmt.Errorf("parse payload: %w", err)
	}

	templateName := common.GetTemplatePath(eventName)
	if _, err := fs.Stat(h.templateFS, "templates/"+templateName); err != nil {
		if h.disableUnknownTemplates {
			return nil, nil
		}
		templateName = common.GetTemplatePath("default")
	}

	return &domain.Notification{
		TemplateName: "github/" + templateName,
		TemplateData: payload,
	}, nil
}
