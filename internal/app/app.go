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
	"github.com/shanth1/hookrelay/internal/service"
	httptransport "github.com/shanth1/hookrelay/internal/transport/http"
)

func Run(ctx, shutdownCtx context.Context, cfg *config.Config) {
	logger := log.FromContext(ctx)

	handlers, err := initInboundHandlers(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize inbound processors")
	}

	notifiers, err := initOutboundAdapters(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize outbound adapters")
	}

	webhookService := service.New(handlers, notifiers, logger)

	runHTTPServer(ctx, shutdownCtx, cfg, webhookService, logger)
	logger.Info().Msg("application shutdown complete")
}

func initInboundHandlers(cfg *config.Config, logger log.Logger) (map[config.WebhookName]ports.WebhookHandler, error) {
	handlers := make(map[config.WebhookName]ports.WebhookHandler)
	for _, webhookCfg := range cfg.Webhooks {
		var handler ports.WebhookHandler
		var err error
		switch webhookCfg.Type {
		case config.WebhookTypeGitHub:
			handler, err = github.NewHandler(webhookCfg.Secret)
			if err != nil {
				return nil, fmt.Errorf("failed to create github processor: %w", err)
			}
		case config.WebhookTypeKanboard:
			handler, err = kanboard.NewHandler(webhookCfg.Secret)
			if err != nil {
				return nil, fmt.Errorf("failed to create kanboard processor: %w", err)
			}
		case config.WebhookTypeCustom:
			handler = custom.NewHandler(webhookCfg.Secret)
		default:
			return nil, fmt.Errorf("unknown webhook handler type '%s' for '%s'", webhookCfg.Type, webhookCfg.Name)
		}
		handlers[config.WebhookName(webhookCfg.Name)] = handler
		logger.Info().Str("name", string(webhookCfg.Name)).Str("type", string(webhookCfg.Type)).Msg("registered webhook handler")
	}

	return handlers, nil

}

func initOutboundAdapters(cfg *config.Config, logger log.Logger) (map[config.NotifierName]ports.Notifier, error) {
	notifiers := make(map[config.NotifierName]ports.Notifier)
	for _, notifierCfg := range cfg.Notifiers {
		var notifier ports.Notifier
		var err error
		switch notifierCfg.Type {
		case config.NotifierTypeTelegram:
			var settings config.TelegramSettings
			if err = notifierCfg.DecodeSettings(&settings); err != nil {
				return nil, fmt.Errorf("failed to decode telegram settings for '%s': %w", notifierCfg.Name, err)
			}
			notifier = telegram.NewSender(settings)
		case config.NotifierTypeEmail:
			var settings config.EmailSettings
			if err = notifierCfg.DecodeSettings(&settings); err != nil {
				return nil, fmt.Errorf("failed to decode email settings for '%s': %w", notifierCfg.Name, err)
			}
			notifier = email.NewSender(settings)
		default:
			return nil, fmt.Errorf("unknown sender type '%s' for '%s'", notifierCfg.Type, notifierCfg.Name)
		}
		notifiers[notifierCfg.Name] = notifier
		logger.Info().Str("name", string(notifierCfg.Name)).Str("type", string(notifierCfg.Type)).Msg("registered notifier")
	}
	return notifiers, nil
}

func runHTTPServer(
	ctx, shutdownCtx context.Context,
	cfg *config.Config,
	service ports.Service,
	logger log.Logger,
) {
	api := httptransport.NewAPI(service, logger, cfg)

	httpHandler := httptransport.NewRouter(api, logger)

	server := httptransport.NewServer(cfg.Addr, httpHandler)

	go func() {
		logger.Info().Msgf("starting HTTP server on %s", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("http server failed")
		}
	}()

	<-ctx.Done()
	logger.Info().Msg("shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("http server graceful shutdown failed")
	}
}
