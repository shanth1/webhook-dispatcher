package app

import (
	"context"
	"net/http"

	"github.com/shanth1/gitrelay/internal/config"
	"github.com/shanth1/gitrelay/internal/handler"
	"github.com/shanth1/gitrelay/internal/service"
	"github.com/shanth1/gitrelay/internal/templates"
	"github.com/shanth1/gotools/log"
)

func Run(ctx, shutdownCtx context.Context, cfg *config.Config) {
	logger := log.FromContext(ctx)

	templates, err := templates.LoadTemplates()
	if err != nil {
		logger.Fatal().Err(err).Msg("load templates")
	}

	notifier, err := service.NewNotifier(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("initialize notifier service")
	}

	handler := handler.New(cfg, templates, notifier, logger)

	mux := http.NewServeMux()
	mux.Handle("/webhook", handler)

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}

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
