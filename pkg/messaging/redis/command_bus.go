package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/messaging"
	"github.com/mikaelchan/hamster/pkg/serializer"
)

var _ messaging.CommandBus = &CommandBus{}

type CommandBus struct {
	*Bus
	handlers map[domain.Type]domain.CommandHandler
	pubSub   *redis.PubSub
}

func NewCommandBus(ctx context.Context, cfg Config, factory *serializer.Factory) messaging.CommandBus {
	cb := &CommandBus{
		Bus:      NewBus(cfg, factory),
		handlers: make(map[domain.Type]domain.CommandHandler),
	}
	cb.pubSub = cb.client.Subscribe(ctx)

	go cb.listen(ctx)
	return cb
}

func (b *CommandBus) Register(ctx context.Context, contract domain.Type, handler domain.CommandHandler) error {
	b.Lock()
	b.handlers[contract] = handler
	b.Unlock()
	return b.pubSub.Subscribe(ctx, string(contract))
}

func (b *CommandBus) Dispatch(ctx context.Context, cmd domain.Command) error {
	return b.publish(ctx, cmd.Type(), cmd)
}

func (b *CommandBus) listen(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case msg := <-b.pubSub.Channel():
		b.RLock()
		listener, exists := b.handlers[domain.Type(msg.Channel)]
		b.RUnlock()
		if exists {
			processFunc := func(ctx context.Context) error {
				cmd, err := b.factory.DeserializeNew(domain.Type(msg.Channel), []byte(msg.Payload))
				if err != nil {
					return fmt.Errorf("unmarshal command: %w", err)
				}

				return listener(ctx, cmd.(domain.Command))
			}

			go b.processMessage(ctx, processFunc)
		}
	}
}

func (b *CommandBus) Close(ctx context.Context) error {
	if err := b.pubSub.Close(); err != nil {
		return err
	}
	return b.close()
}
