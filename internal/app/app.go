package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/adapters/inbound/custom"
	"github.com/shanth1/hookrelay/internal/adapters/inbound/github"
	"github.com/shanth1/hookrelay/internal/adapters/inbound/kanboard"
	"github.com/shanth1/hookrelay/internal/adapters/outbound/email"
	"github.com/shanth1/hookrelay/internal/adapters/outbound/telegram"
	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/core/ports"
	"github.com/shanth1/hookrelay/internal/service/webhook"
	httptransport "github.com/shanth1/hookrelay/internal/transport/http"
	"github.com/shanth1/hookrelay/internal/transport/http/router"
)

func Run(ctx, shutdownCtx context.Context, cfg *config.Config) {
	logger := log.FromContext(ctx)

	notifiers, err := initOutboundAdapters(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize outbound adapters")
	}

	processors, err := initInboundProcessors()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize inbound processors")
	}

	webhookService := webhook.NewService(notifiers, processors, logger)

	runHTTPServer(ctx, shutdownCtx, cfg, webhookService, processors, logger)
	logger.Info().Msg("application shutdown complete")
}

func runHTTPServer(
	ctx, shutdownCtx context.Context,
	cfg *config.Config,
	service ports.WebhookService,
	processors map[config.WebhookType]webhook.Processor,
	logger log.Logger,
) {
	router := router.New(cfg, service, processors, logger)
	httpHandler := router.Handler()
	server := httptransport.NewServer(cfg.Addr, httpHandler)

	go func() {
		logger.Info().Msgf("starting HTTP server on %s", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("http server failed")
		}
	}()

	<-ctx.Done()
	logger.Info().Msg("shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("http server graceful shutdown failed")
	}
}

func initOutboundAdapters(cfg *config.Config, logger log.Logger) (map[string]ports.Notifier, error) {
	notifiers := make(map[string]ports.Notifier)
	for _, senderCfg := range cfg.Senders {
		var s ports.Notifier
		var err error
		switch senderCfg.Type {
		case config.SenderTypeTelegram:
			var settings config.TelegramSettings
			if err = senderCfg.DecodeSenderSettings(&settings); err != nil {
				return nil, fmt.Errorf("failed to decode telegram settings for '%s': %w", senderCfg.Name, err)
			}
			s = telegram.NewSender(settings)
		case config.SenderTypeEmail:
			var settings config.EmailSettings
			if err = senderCfg.DecodeSenderSettings(&settings); err != nil {
				return nil, fmt.Errorf("failed to decode email settings for '%s': %w", senderCfg.Name, err)
			}
			s = email.NewSender(settings)
		default:
			return nil, fmt.Errorf("unknown sender type '%s' for '%s'", senderCfg.Type, senderCfg.Name)
		}
		notifiers[senderCfg.Name] = s
		logger.Info().Str("name", senderCfg.Name).Str("type", string(senderCfg.Type)).Msg("registered notifier")
	}
	return notifiers, nil
}

func initInboundProcessors() (map[config.WebhookType]webhook.Processor, error) {
	githubProcessor, err := github.NewProcessor()
	if err != nil {
		return nil, fmt.Errorf("failed to create github processor: %w", err)
	}
	kanboardProcessor, err := kanboard.NewProcessor()
	if err != nil {
		return nil, fmt.Errorf("failed to create kanboard processor: %w", err)
	}
	customProcessor := custom.NewProcessor()

	return map[config.WebhookType]webhook.Processor{
		config.WebhookTypeGitHub:   githubProcessor,
		config.WebhookTypeKanboard: kanboardProcessor,
		config.WebhookTypeCustom:   customProcessor,
	}, nil
}
