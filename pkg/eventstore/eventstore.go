package eventstore

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
)

type EventIterator interface {
	HasNext(ctx context.Context) bool
	Next(ctx context.Context) (domain.Event, error)
	Close(ctx context.Context) error
}

type EventData interface {
	StreamID() string
	Type() string
	Payload() []byte
}

// EventStore defines the methods for storing and retrieving events
type EventStore interface {
	Append(ctx context.Context, streamID domain.ID, events ...domain.Event) error
	Load(ctx context.Context, streamID domain.ID, from domain.Version) (EventIterator, error)
}
