package service

import "context"

type Sender interface {
	Send(ctx context.Context, target string, message string) error
}
