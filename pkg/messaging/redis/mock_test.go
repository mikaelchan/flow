package redis_test

import (
	"github.com/mikaelchan/hamster/pkg/domain"
)

// MockEvent is a simple implementation of the Event interface for testing
type MockEvent struct {
	domain.BaseEvent
}

func (e *MockEvent) Type() domain.Type {
	return "mock.event"
}

// MockAggregateRoot is a simple implementation of the AggregateRoot interface for testing
type MockAggregateRoot struct {
	domain.BaseAggregateRoot
}

func (m *MockAggregateRoot) Type() domain.Type {
	return "mock"
}

func (m *MockAggregateRoot) Apply(event domain.Event) error {
	return nil
}
