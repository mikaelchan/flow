package mongo_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mikaelchan/hamster/pkg/domain"
	mongoes "github.com/mikaelchan/hamster/pkg/eventstore/mongo"
	"github.com/mikaelchan/hamster/pkg/serializer"
	"github.com/mikaelchan/hamster/pkg/serializer/json"

	"go.mongodb.org/mongo-driver/bson"
)

type MockAggregateRoot struct {
	domain.BaseAggregateRoot
}

func (m *MockAggregateRoot) Type() domain.Type {
	return "mock"
}

func (m *MockAggregateRoot) Apply(event domain.Event) error {
	return nil
}

type MockEvent struct {
	domain.BaseEvent
}

var db *mongo.Database

func (e *MockEvent) Type() domain.Type {
	return "mock.event"
}

func TestMain(m *testing.M) {
	json.RegisterJSON(&MockEvent{})
	clientOptions := options.Client().ApplyURI("mongodb://nas.myhome:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	db = client.Database("test")

	m.Run()
	db.Collection("events").DeleteMany(context.Background(), bson.M{})
	client.Disconnect(context.Background())
}

func TestEventStore_Append(t *testing.T) {
	factory := serializer.GetFactory()
	store, err := mongoes.NewEventStore(context.Background(), db, "events", factory)
	if err != nil {
		t.Fatalf("Failed to create event store: %v", err)
	}
	root := &MockAggregateRoot{}

	event := &MockEvent{domain.NewBaseEvent("1234", root)}
	err = store.Append(context.Background(), []domain.Event{event})
	if err != nil {
		t.Fatalf("Failed to append event: %v", err)
	}

	var result bson.M
	err = db.Collection("events").FindOne(context.Background(), bson.M{"event_id": event.ID().String()}).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to find appended event: %v", err)
	}
}

func TestEventStore_Load(t *testing.T) {
	factory := serializer.GetFactory()
	store, err := mongoes.NewEventStore(context.Background(), db, "events", factory)
	if err != nil {
		t.Fatalf("Failed to create event store: %v", err)
	}
	event := &MockEvent{domain.NewBaseEvent("1234", &MockAggregateRoot{})}
	pyload, err := factory.Serialize(event)
	if err != nil {
		t.Fatalf("Failed to serialize event: %v", err)
	}

	_, err = db.Collection("events").InsertOne(context.Background(), bson.M{
		"stream_id":  event.StreamID().String(),
		"event_id":   event.ID().String(),
		"event_type": event.Type(),
		"version":    event.StreamVersion(),
		"payload":    pyload,
	})
	if err != nil {
		t.Fatalf("Failed to insert event: %v", err)
	}

	iter, err := store.Load(context.Background(), event.StreamID(), 1)
	if err != nil {
		t.Fatalf("Failed to load events: %v", err)
	}
	defer iter.Close(context.Background())

	if !iter.HasNext(context.Background()) {
		t.Fatalf("Expected to have next event")
	}

	loadedEvent, err := iter.Next(context.Background())
	if err != nil {
		t.Fatalf("Failed to get next event: %v", err)
	}

	if loadedEvent.ID() != event.ID() {
		t.Errorf("Expected event ID %s, got %s", event.ID(), loadedEvent.ID())
	}
}
