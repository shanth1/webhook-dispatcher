package server

import (
	"fmt"
	"net/http"

	"github.com/shanth1/gitrelay/internal/config"
	"github.com/shanth1/gitrelay/internal/middleware"
	"github.com/shanth1/gitrelay/internal/processor"
	"github.com/shanth1/gitrelay/internal/verifier"
)

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()

	githubProcessor := processor.NewGithubProcessor(s.templates)
	customProcessor := processor.NewCustomProcessor()

	recipientMap := make(map[string]config.Recipient, len(s.cfg.Recipients))
	for _, r := range s.cfg.Recipients {
		recipientMap[r.Name] = r
	}

	for _, hookCfg := range s.cfg.Webhooks {
		var v verifier.Verifier
		var p processor.WebhookProcessor

		switch hookCfg.Type {
		case config.WebhookTypeGitHub:
			v = &verifier.GithubVerifier{Secret: hookCfg.Secret}
			p = githubProcessor
		case config.WebhookTypeCustom:
			v = &verifier.CustomVerifier{Secret: hookCfg.Secret}
			p = customProcessor
		default:
			s.logger.Error().Str("type", string(hookCfg.Type)).Msg("unknown webhook type in config, skipping")
			continue
		}

		resolvedRecipients := make([]config.Recipient, 0, len(hookCfg.Recipients))
		for _, recipientName := range hookCfg.Recipients {
			if recipient, ok := recipientMap[recipientName]; ok {
				resolvedRecipients = append(resolvedRecipients, recipient)
			} else {
				s.logger.Warn().
					Str("webhook_name", hookCfg.Name).
					Str("recipient_name", recipientName).
					Msg("recipient name not found in global definitions, skipping for this webhook")
			}
		}

		if len(resolvedRecipients) == 0 {
			s.logger.Warn().Str("webhook_name", hookCfg.Name).Msg("webhook has no valid recipients configured, it will not send any notifications")
		}

		s.logger.Info().
			Str("name", hookCfg.Name).
			Str("path", hookCfg.Path).
			Str("type", string(hookCfg.Type)).
			Int("recipient_count", len(resolvedRecipients)).
			Msg("registering webhook handler")

		handlerChain := middleware.Chain(
			s.handleWebhook(p, resolvedRecipients),
			middleware.WithLogger(s.logger),
			middleware.WithMethod(http.MethodPost),
			middleware.WithVerification(v),
		)
		mux.Handle(fmt.Sprintf("POST %s", hookCfg.Path), handlerChain)
	}

	mux.Handle("GET /health", middleware.Chain(
		http.HandlerFunc(s.handleHealthCheck),
		middleware.WithLogger(s.logger),
	))
	mux.Handle("/", middleware.Chain(
		http.HandlerFunc(s.handleRoot),
		middleware.WithLogger(s.logger),
	))

	return mux
}
