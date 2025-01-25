package domain_test

import (
	"errors"
	"testing"

	"github.com/mikaelchan/hamster/pkg/domain"
)

func TestBaseAggregateRoot_Track(t *testing.T) {
	root := &MockAggregateRoot{}

	event := &MockEvent{BaseEvent: domain.NewBaseEvent("1234", root)}

	err := domain.Track(root, event)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(root.UncommittedEvents()) != 1 {
		t.Errorf("Expected 1 uncommitted event, got %d", len(root.UncommittedEvents()))
	}

	if root.UncommittedVersion() != 1 {
		t.Errorf("Expected uncommitted version 1, got %d", root.UncommittedVersion())
	}
}

func TestBaseAggregateRoot_ClearUncommittedEvents(t *testing.T) {
	root := &MockAggregateRoot{}

	event := &MockEvent{BaseEvent: domain.NewBaseEvent("1234", root)}
	_ = domain.Track(root, event)

	root.ClearUncommittedEvents()

	if len(root.UncommittedEvents()) != 0 {
		t.Errorf("Expected 0 uncommitted events, got %d", len(root.UncommittedEvents()))
	}

	if root.Version() != 1 {
		t.Errorf("Expected version 1, got %d", root.Version())
	}
}

func TestBaseAggregateRoot_SetID(t *testing.T) {
	root := &domain.BaseAggregateRoot{}

	err := root.SetID("1234")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if root.ID() != "1234" {
		t.Errorf("Expected ID %s, got %s", "1234", root.ID())
	}

	// Attempt to set ID again
	err = root.SetID("5678")
	if !errors.Is(err, domain.ErrIDAlreadySet) {
		t.Errorf("Expected error %v, got %v", domain.ErrIDAlreadySet, err)
	}
}
