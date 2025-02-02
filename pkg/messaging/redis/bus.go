package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/serializer"
)

// Bus contains common functionality for Redis-based buses
type Bus struct {
	client        *redis.Client
	handleTimeout time.Duration
	factory       *serializer.Factory
	sync.RWMutex
}

func NewBus(config Config, factory *serializer.Factory) *Bus {
	if config.HandleTimeout == 0 {
		config.HandleTimeout = 30 * time.Second
	}
	return &Bus{
		client:        config.Client,
		handleTimeout: config.HandleTimeout,
		factory:       factory,
	}
}

// processMessage executes a process with timeout
func (b *Bus) processMessage(ctx context.Context, process func(context.Context) error) {
	ctx, cancel := context.WithTimeout(ctx, b.handleTimeout)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- process(ctx)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			fmt.Printf("error handling message: %v\n", err)
		}
	case <-ctx.Done():
		fmt.Printf("message handling timed out after %v\n", b.handleTimeout)
	}
}

// publish sends a message to Redis
func (b *Bus) publish(ctx context.Context, typ domain.Type, msg domain.HasType) error {
	msgBytes, err := b.factory.Serialize(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	err = b.client.Publish(ctx, typ.String(), msgBytes).Err()
	if err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	return nil
}
