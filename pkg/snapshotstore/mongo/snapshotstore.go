package mongo

import (
	"context"
	"errors"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/serializer"
	"github.com/mikaelchan/hamster/pkg/snapshotstore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SnapshotStore implements the SnapshotStore interface using MongoDB
type snapshotStore struct {
	collection *mongo.Collection
	factory    serializer.Factory
	policy     snapshotstore.SnapshotPolicy
}

// NewMongoSnapshotStore creates a new MongoSnapshotStore
func NewMongoSnapshotStore(db *mongo.Database, collectionName string, factory serializer.Factory, policy snapshotstore.SnapshotPolicy) snapshotstore.SnapshotStore {
	return &snapshotStore{
		collection: db.Collection(collectionName),
		factory:    factory,
		policy:     policy,
	}
}

// Save stores an aggregate's snapshot in the MongoDB collection
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

	filter := bson.M{"stream_id": snapshot.StreamID(), "type": snapshot.Type()}
	update := bson.M{
		"$set": bson.M{
			"version": snapshot.Version(),
			"state":   snapshot.State(),
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err = s.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// Load retrieves an aggregate's snapshot from the MongoDB collection
func (s *snapshotStore) Load(ctx context.Context, streamID domain.ID, root domain.AggregateRoot) error {
	var mongoSnapshot Snapshot
	filter := bson.M{"stream_id": streamID.String()}
	err := s.collection.FindOne(ctx, filter).Decode(&mongoSnapshot)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("snapshot not found")
	}
	if err != nil {
		return err
	}
	// check if the type is the same
	if root.Type() != domain.Type(mongoSnapshot.Type()) {
		return errors.New("snapshot type mismatch")
	}

	err = s.factory.Deserialize(mongoSnapshot.StateField, root)
	if err != nil {
		return err
	}
	return nil
}
