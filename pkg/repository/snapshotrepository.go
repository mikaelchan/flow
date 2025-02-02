package repository

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/eventstore"
	"github.com/mikaelchan/hamster/pkg/messaging"
	"github.com/mikaelchan/hamster/pkg/snapshotstore"
)

var _ Repository = &snapshotRepository{}

type snapshotRepository struct {
	*eventRepository
	store snapshotstore.SnapshotStore
}

func NewSnapshotRepository(store eventstore.EventStore, snapshotStore snapshotstore.SnapshotStore, bus messaging.EventBus) Repository {
	return &snapshotRepository{
		eventRepository: NewEventRepository(store, bus),
		store:           snapshotStore,
	}
}

func (r *snapshotRepository) Load(ctx context.Context, id domain.ID, root domain.AggregateRoot) error {
	err := r.store.Load(ctx, id, root)
	if err != nil {
		// TODO: add log
		// fallback to event repository
		return r.eventRepository.Load(ctx, id, root)
	}
	err = r.eventRepository.LoadFrom(ctx, id, root, root.UncommittedVersion())
	if err != nil {
		return err
	}
	return nil
}

func (r *snapshotRepository) Save(ctx context.Context, root domain.AggregateRoot) error {
	go func() {
		err := r.store.Save(ctx, root)
		if err != nil {
			// TODO: add log
		}
	}()

	return r.eventRepository.Save(ctx, root)
}
