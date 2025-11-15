package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/core/domain"
	"github.com/shanth1/hookrelay/internal/core/ports"
)

type Service struct {
	handlers  map[config.WebhookName]ports.WebhookHandler
	notifiers map[config.NotifierName]ports.Notifier
	logger    log.Logger
}

var _ ports.Service = (*Service)(nil)

func New(
	handlers map[config.WebhookName]ports.WebhookHandler,
	notifiers map[config.NotifierName]ports.Notifier,
	logger log.Logger,
) ports.Service {
	return &Service{
		logger:    logger,
		handlers:  handlers,
		notifiers: notifiers,
	}
}

func (s *Service) ProcessWebhook(ctx context.Context, webhookName config.WebhookName, req ports.WebhookRequest, recipients []config.Recipient) error {
	webhookHandler, ok := s.handlers[webhookName]
	if !ok {
		return fmt.Errorf("no handler registered for webhook name: %s", webhookHandler)
	}

	notification, err := webhookHandler.Handle(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to process request payload: %w", err)
	}

	if notification == nil {
		s.logger.Info().Str("name", string(webhookName)).Msg("handler returned no notification, skipping broadcast")
		return nil
	}

	s.broadcast(ctx, recipients, *notification)
	return nil
}

func (s *Service) broadcast(ctx context.Context, recipients []config.Recipient, notification domain.Notification) {
	var wg sync.WaitGroup
	logger := log.FromContext(ctx)

	for _, recipient := range recipients {
		notifier, ok := s.notifiers[recipient.Notifier]
		if !ok {
			logger.Error().Str("recipient", recipient.Name).Str("notifier", string(recipient.Notifier)).Msg("notifier not found")
			continue
		}

		wg.Add(1)
		go func(rcp config.Recipient, ntf ports.Notifier) {
			defer wg.Done()
			if err := ntf.Send(ctx, rcp.Target, notification); err != nil {
				logger.Error().
					Str("recipient", rcp.Name).
					Str("notifier_name", string(rcp.Notifier)).
					Err(err).
					Msg("failed to send notification")
			}
		}(recipient, notifier)
	}
	wg.Wait()
}
