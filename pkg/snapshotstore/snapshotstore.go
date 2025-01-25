package snapshotstore

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
)

// Snapshot represents a snapshot of an aggregate's state.
type Snapshot interface {
	StreamID() string
	Type() string
	Version() uint64
	State() []byte
}

type ShouldSnapshotFunc func(ctx context.Context, aggregate domain.AggregateRoot) bool

type SnapshotPolicy struct {
	ShouldSnapshot ShouldSnapshotFunc
}

type SnapshotStore interface {
	Save(ctx context.Context, aggregate domain.AggregateRoot) error
	Load(ctx context.Context, streamID domain.ID, root domain.AggregateRoot) error
}
