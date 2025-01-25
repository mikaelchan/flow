package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/messaging"
	"github.com/mikaelchan/hamster/pkg/serializer"
)

var _ messaging.EventBus = &EventBus{}

type EventBus struct {
	*Bus
	subscribers map[domain.Type][]domain.EventListener
	pubSub      *redis.PubSub
}

func NewEventBus(ctx context.Context, cfg Config, factory *serializer.Factory) messaging.EventBus {
	eb := &EventBus{
		Bus:         NewBus(cfg, factory),
		subscribers: make(map[domain.Type][]domain.EventListener),
	}
	eb.pubSub = eb.client.Subscribe(ctx)

	go eb.listen(ctx)
	return eb
}

// Subscribe registers a handler for an event topic
func (b *EventBus) Subscribe(ctx context.Context, topic domain.Type, listener domain.EventListener) error {
	b.Lock()
	b.subscribers[topic] = append(b.subscribers[topic], listener)
	b.Unlock()
	return b.pubSub.Subscribe(ctx, topic.String())
}

// Publish publishes an event to Redis
func (b *EventBus) Publish(ctx context.Context, event domain.Event) error {
	return b.publish(ctx, event.Type(), event)
}

// listen listens for events on all subscribed channels
func (b *EventBus) listen(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case msg := <-b.pubSub.Channel():
		b.RLock()
		handlers, exists := b.subscribers[domain.Type(msg.Channel)]
		b.RUnlock()
		if exists {
			for _, handler := range handlers {
				processFunc := func(ctx context.Context) error {
					event, err := b.factory.DeserializeNew(domain.Type(msg.Channel), []byte(msg.Payload))
					if err != nil {
						return fmt.Errorf("unmarshal event: %w", err)
					}
					return handler(ctx, event.(domain.Event))
				}
				go b.processMessage(ctx, processFunc)
			}
		}
	}
}

// Close closes the Redis connection
func (b *EventBus) Close(ctx context.Context) error {
	if err := b.pubSub.Close(); err != nil {
		return err
	}
	return b.client.Close()
}
