package notifier

import (
	"context"
	"fmt"
	"sync"

	"github.com/shanth1/gitrelay/internal/config"
	"github.com/shanth1/gotools/log"
)

type Notifier struct {
	senders map[string]Sender
	logger  log.Logger
}

func NewNotifier(cfg *config.Config, logger log.Logger) (*Notifier, error) {
	senders := make(map[string]Sender)
	for _, senderCfg := range cfg.Senders {
		if _, exists := senders[senderCfg.Name]; exists {
			return nil, fmt.Errorf("sender with name '%s' is already registered", senderCfg.Name)
		}

		var s Sender
		var err error

		switch senderCfg.Type {
		case config.SenderTypeTelegram:
			var settings config.TelegramSettings
			if err = senderCfg.DecodeSenderSettings(&settings); err != nil {
				return nil, fmt.Errorf("failed to decode telegram settings for sender '%s': %w", senderCfg.Name, err)
			}
			s = NewTelegramSender(settings)
		case config.SenderTypeEmail:
			var settings config.EmailSettings
			if err = senderCfg.DecodeSenderSettings(&settings); err != nil {
				return nil, fmt.Errorf("failed to decode email settings for sender '%s': %w", senderCfg.Name, err)
			}
			s = NewEmailSender(settings)
		default:
			return nil, fmt.Errorf("unknown sender type '%s' for sender '%s'", senderCfg.Type, senderCfg.Name)
		}

		senders[senderCfg.Name] = s
		logger.Info().Str("name", senderCfg.Name).Str("type", string(senderCfg.Type)).Msg("registered sender")
	}

	return &Notifier{senders: senders, logger: logger}, nil
}

func (n *Notifier) Broadcast(ctx context.Context, recipients []config.Recipient, message string) {
	var wg sync.WaitGroup
	for _, r := range recipients {
		sender, ok := n.senders[r.Sender]
		if !ok {
			n.logger.Error().Str("recipient", r.Name).Str("sender", r.Sender).Msg("sender not found for recipient")
			continue
		}

		wg.Add(1)
		go func(recipient config.Recipient, s Sender) {
			defer wg.Done()
			if err := s.Send(ctx, recipient.Target, message); err != nil {
				n.logger.Error().
					Str("recipient", r.Name).
					Str("sender_name", recipient.Sender).
					Str("target", recipient.Target).
					Err(err).
					Msg("failed to send notification")
			}
		}(r, sender)
	}
	wg.Wait()
}
