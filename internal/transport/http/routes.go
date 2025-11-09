package httptransport

import (
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/transport/http/middleware"
)

func NewRouter(api *API, logger log.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", api.handleHealthCheck)
	mux.HandleFunc("GET /webhooks", api.handleListWebhooks)
	mux.HandleFunc("GET /", api.handleRoot)

	for _, hookCfg := range api.config.Webhooks {
		handler := api.webhookHandlerFactory(hookCfg)

		postOnlyHandler := middleware.WithMethod(http.MethodPost)(handler)

		logger.Info().Str("path", hookCfg.Path).Str("type", string(hookCfg.Type)).Msg("registering webhook handler")
		mux.Handle("POST "+hookCfg.Path, postOnlyHandler)
	}

	return middleware.Chain(
		mux,
		middleware.WithRecovery(logger),
		middleware.WithLogger(logger),
	)
}
