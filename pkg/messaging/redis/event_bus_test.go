package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/mikaelchan/hamster/pkg/domain"
	redismessaging "github.com/mikaelchan/hamster/pkg/messaging/redis"
	"github.com/mikaelchan/hamster/pkg/serializer/json"
)

type MockEventListener struct {
	receivedEvents []*MockEvent
}

func (l *MockEventListener) Handle(ctx context.Context, event domain.Event) error {
	l.receivedEvents = append(l.receivedEvents, event.(*MockEvent))
	return nil
}

func TestMain(m *testing.M) {
	json.RegisterJSON(&MockEvent{})
	m.Run()
}

func TestRedisEventBus_SubscribeAndPublish(t *testing.T) {
	ctx := context.Background()
	// Create a new RedisEventBus
	bus := redismessaging.NewEventBus(ctx,
		redismessaging.Config{
			Client:        redis.NewClient(&redis.Options{Addr: "nas.myhome:6379"}),
			HandleTimeout: 5 * time.Second,
		},
	)
	defer bus.Close(ctx)

	// Create a mock event listener
	listener := &MockEventListener{}

	// Subscribe to an event type
	eventType := domain.Type("mock.event")
	err := bus.Subscribe(ctx, eventType, listener.Handle)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Create a mock event
	mockEvent := &MockEvent{
		BaseEvent: domain.NewBaseEvent("1234", &MockAggregateRoot{}),
	}

	// Publish the event
	err = bus.Publish(ctx, mockEvent)
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	// Allow some time for the event to be processed
	time.Sleep(1 * time.Second)

	// Verify the event was received
	if len(listener.receivedEvents) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(listener.receivedEvents))
	}

	if listener.receivedEvents[0].ID() != mockEvent.ID() {
		t.Errorf("Expected event ID %s, got %s", mockEvent.ID(), listener.receivedEvents[0].ID())
	}
}
