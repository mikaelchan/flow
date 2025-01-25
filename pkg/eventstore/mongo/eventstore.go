package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/eventstore"
	"github.com/mikaelchan/hamster/pkg/serializer"
)

var _ eventstore.EventStore = &EventStore{}

// EventStore implements the EventStore interface using MongoDB
type EventStore struct {
	collection *mongo.Collection
	factory    *serializer.Factory
}

// NewEventStore creates a new MongoEventStore and initializes indexes
func NewEventStore(ctx context.Context, db *mongo.Database, collectionName string, factory *serializer.Factory) (*EventStore, error) {
	store := &EventStore{
		collection: db.Collection(collectionName),
		factory:    factory,
	}

	// Define indexes
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "stream_id", Value: 1},
				{Key: "version", Value: 1},
			},
			Options: options.Index().SetName("stream_id_version"),
		},
		{
			Keys: bson.D{
				{Key: "event_id", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetName("unique_event_id"),
		},
	}

	// Create indexes
	for _, indexModel := range indexModels {
		_, err := store.collection.Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			return nil, err // Return error if index creation fails
		}
	}

	return store, nil
}

// Append adds events to the MongoDB collection
func (s *EventStore) Append(ctx context.Context, events []domain.Event) error {
	if len(events) == 0 {
		return errors.New("no events to append")
	}

	var documents []interface{}
	for _, event := range events {
		payload, err := s.factory.Serialize(event)
		if err != nil {
			return err
		}

		documents = append(documents, bson.M{
			"stream_id":  event.StreamID().String(),
			"event_id":   event.ID().String(),
			"event_type": event.Type(),
			"version":    event.StreamVersion(),
			"payload":    payload,
		})
	}

	_, err := s.collection.InsertMany(ctx, documents)
	return err
}

// Load retrieves events from the MongoDB collection, starting from the given version(inclusive)
func (s *EventStore) Load(ctx context.Context, streamID domain.ID, from domain.Version) (eventstore.EventIterator, error) {
	filter := bson.M{"stream_id": streamID.String(), "version": bson.M{"$gte": from}}
	findOptions := options.Find().SetSort(bson.D{{Key: "version", Value: 1}})
	cursor, err := s.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return &eventstore.EmptyIterator{}, err
	}

	return &eventIterator{cursor: cursor, factory: s.factory}, nil
}
