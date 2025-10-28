package app

import (
	"context"
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/notifier"
	"github.com/shanth1/hookrelay/internal/server"
	"github.com/shanth1/hookrelay/internal/templates"
)

func Run(ctx, shutdownCtx context.Context, cfg *config.Config) {
	logger := log.FromContext(ctx)

	templates, err := templates.LoadTemplates()
	if err != nil {
		logger.Fatal().Err(err).Msg("load templates")
	}

	notifier, err := notifier.NewNotifier(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("initialize notifier service")
	}

	server := server.New(cfg, templates, notifier, logger)

	go func() {
		logger.Info().Msgf("staring server on %s", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("start server")
		}
	}()

	<-ctx.Done()

	logger.Info().Msg("shutting down server...")

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("shutdown server")
	}

	logger.Info().Msg("server shutdown completed")
}
