package repository

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/eventstore"
	"github.com/mikaelchan/hamster/pkg/messaging"
)

type eventRepository struct {
	store eventstore.EventStore
	bus   messaging.EventBus
}

func NewEventRepository(store eventstore.EventStore, bus messaging.EventBus) *eventRepository {
	return &eventRepository{
		store: store,
		bus:   bus,
	}
}

func (r *eventRepository) Save(ctx context.Context, root domain.AggregateRoot) error {
	// Collect uncommitted events from the aggregate root
	events := root.UncommittedEvents()

	// Save these events to the event store
	if err := r.store.Append(ctx, events); err != nil {
		return err
	}

	// Publish each event to the event bus
	for _, e := range events {
		if err := r.bus.Publish(ctx, e); err != nil {
			return err
		}
	}

	// Clear uncommitted events to indicate they're now persisted
	root.ClearUncommittedEvents()
	return nil
}

func (r *eventRepository) Load(ctx context.Context, id domain.ID, root domain.AggregateRoot) error {
	return r.LoadFrom(ctx, id, root, 0)
}

func (r *eventRepository) LoadFrom(ctx context.Context, id domain.ID, root domain.AggregateRoot, from domain.Version) error {
	it, err := r.store.Load(ctx, id, from)
	if err != nil {
		return err
	}
	defer it.Close(ctx)

	for it.HasNext(ctx) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			event, err := it.Next(ctx)
			if err != nil {
				return err
			}
			root.Apply(event)
		}

	}
	return nil
}
