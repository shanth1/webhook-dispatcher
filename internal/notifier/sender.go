package notifier

import "context"

type Sender interface {
	Send(ctx context.Context, target string, message string) error
}
