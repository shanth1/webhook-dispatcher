package webhook

import (
	"context"
	"fmt"
	"sync"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/core/domain"
	"github.com/shanth1/hookrelay/internal/core/ports"
)

// Processor - это контракт, который сервис определяет для своих входящих адаптеров.
// Он описывает, как адаптировать сырые данные в доменную модель.
type Processor interface {
	Process(payload []byte, headers map[string]string) (*domain.Notification, error)
}

// Service является реализацией порта WebhookService.
type Service struct {
	notifiers  map[string]ports.Notifier
	processors map[config.WebhookType]Processor
	logger     log.Logger
}

var _ ports.WebhookService = (*Service)(nil)

func NewService(
	notifiers map[string]ports.Notifier,
	processors map[config.WebhookType]Processor,
	logger log.Logger,
) ports.WebhookService {
	return &Service{
		notifiers:  notifiers,
		processors: processors,
		logger:     logger,
	}
}

func (s *Service) ProcessWebhook(ctx context.Context, req ports.WebhookRequest, recipients []config.Recipient) error {
	processor, ok := s.processors[req.WebhookType]
	if !ok {
		return fmt.Errorf("no processor registered for webhook type: %s", req.WebhookType)
	}

	notification, err := processor.Process(req.Payload, req.Headers)
	if err != nil {
		return fmt.Errorf("failed to process request payload: %w", err)
	}

	s.broadcast(ctx, recipients, *notification)
	return nil
}

func (s *Service) broadcast(ctx context.Context, recipients []config.Recipient, notification domain.Notification) {
	var wg sync.WaitGroup
	logger := log.FromContext(ctx)

	for _, recipient := range recipients {
		notifier, ok := s.notifiers[recipient.Sender]
		if !ok {
			logger.Error().Str("recipient", recipient.Name).Str("sender", recipient.Sender).Msg("sender not found")
			continue
		}

		wg.Add(1)
		go func(rcp config.Recipient, ntf ports.Notifier) {
			defer wg.Done()
			if err := ntf.Send(ctx, rcp.Target, notification); err != nil {
				logger.Error().
					Str("recipient", rcp.Name).
					Str("sender_name", rcp.Sender).
					Err(err).
					Msg("failed to send notification")
			}
		}(recipient, notifier)
	}
	wg.Wait()
}
