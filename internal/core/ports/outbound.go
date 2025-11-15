package ports

import (
	"context"

	"github.com/shanth1/hookrelay/internal/core/domain"
)

type Notifier interface {
	Send(ctx context.Context, target string, notification domain.Notification) error
}
