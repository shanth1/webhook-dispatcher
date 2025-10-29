package processor

import "net/http"

type WebhookProcessor interface {
	Process(r *http.Request) (message string, err error)
}
