package custom

import (
	"fmt"

	"github.com/shanth1/hookrelay/internal/core/domain"
	"github.com/shanth1/hookrelay/internal/service/webhook"
)

type Processor struct{}

var _ webhook.Processor = (*Processor)(nil)

func NewProcessor() webhook.Processor {
	return &Processor{}
}

func (p *Processor) Process(payload []byte, headers map[string]string) (*domain.Notification, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("request payload is empty")
	}
	return &domain.Notification{Body: string(payload)}, nil
}
