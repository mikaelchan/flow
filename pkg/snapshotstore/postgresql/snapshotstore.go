package postgresql

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/serializer"
	"github.com/mikaelchan/hamster/pkg/snapshotstore"
	"gorm.io/gorm"
)

type snapshotStore struct {
	db      *gorm.DB
	factory *serializer.Factory
	policy  snapshotstore.SnapshotPolicy
}

func NewPostgresSnapshotStore(db *gorm.DB, factory *serializer.Factory, policy snapshotstore.SnapshotPolicy) snapshotstore.SnapshotStore {
	return &snapshotStore{
		db:      db,
		factory: factory,
		policy:  policy,
	}
}

func (s *snapshotStore) Save(ctx context.Context, aggregate domain.AggregateRoot) error {
	if !s.policy.ShouldSnapshot(ctx, aggregate) {
		return nil
	}

	state, err := s.factory.Serialize(aggregate)
	if err != nil {
		return err
	}
	snapshot := Snapshot{
		StreamIDField: aggregate.ID().String(),
		TypeField:     string(aggregate.Type()),
		VersionField:  uint64(aggregate.Version()),
		StateField:    state,
	}

	return s.db.Create(&snapshot).Error
}

func (s *snapshotStore) Load(ctx context.Context, streamID domain.ID, root domain.AggregateRoot) error {
	var snapshot Snapshot
	if err := s.db.Where("stream_id = ?", streamID.String()).First(&snapshot).Error; err != nil {
		return err
	}

	return s.factory.Deserialize(snapshot.StateField, root)
}
