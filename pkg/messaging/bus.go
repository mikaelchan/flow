package messaging

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
)

type CommandBus interface {
	Register(ctx context.Context, contract domain.Type, handler domain.CommandHandler) error
	Dispatch(ctx context.Context, cmd domain.Command) error
	Close(ctx context.Context) error
}

type EventBus interface {
	Subscribe(ctx context.Context, topic domain.Type, listener domain.EventListener) error
	Publish(ctx context.Context, event domain.Event) error
	Close(ctx context.Context) error
}
