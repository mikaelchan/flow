package postgresql_test

import (
	"context"
	"testing"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/eventstore/postgresql"
	"github.com/mikaelchan/hamster/pkg/serializer"
	"github.com/mikaelchan/hamster/pkg/serializer/json"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type MockAggregateRoot struct {
	domain.BaseAggregateRoot
}

func (m *MockAggregateRoot) Type() domain.Type {
	return "mock.aggregate"
}

type MockEvent struct {
	domain.BaseEvent
}

func (m *MockEvent) Type() domain.Type {
	return "mock.event"
}

func TestMain(m *testing.M) {
	json.RegisterJSON(&MockEvent{})
	dsn := "host=nas.myhome user=user password=password dbname=hamster port=5433 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestEventStore_AppendAndLoad(t *testing.T) {
	es := postgresql.NewEventStore(db, serializer.GetFactory())
	t.Run("append", func(t *testing.T) {
		root := &MockAggregateRoot{}
		root.SetID("test")
		event := &MockEvent{domain.NewBaseEvent("1234", root)}
		err := es.Append(context.Background(), event.StreamID(), event)
		if err != nil {
			t.Fatalf("Failed to append event: %v", err)
		}
	})
	t.Run("load", func(t *testing.T) {
		events, err := es.Load(context.Background(), "test", 0)
		if err != nil {
			t.Fatalf("Failed to load events: %v", err)
		}
		if !events.HasNext(context.Background()) {
			t.Fatalf("Expected 1 event, got 0")
		}
		event, err := events.Next(context.Background())
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}
		if event == nil {
			t.Fatalf("Expected 1 event, got nil")
		}
		if event.ID().String() != "1234" {
			t.Fatalf("Expected event ID 1234, got %s", event.ID().String())
		}
	})
	t.Cleanup(func() {
		db.Delete(&postgresql.EventData{}, "stream_id = ?", "test")
	})
}
