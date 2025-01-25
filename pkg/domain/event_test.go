package domain_test

import (
	"testing"

	"github.com/mikaelchan/hamster/pkg/domain"
)

func TestBaseEvent(t *testing.T) {
	mockRoot := &MockAggregateRoot{}
	mockRoot.SetID("5678")
	baseEvent := domain.NewBaseEvent("1234", mockRoot)

	if baseEvent.ID() != "1234" {
		t.Errorf("Expected ID %s, got %s", "1234", baseEvent.ID())
	}

	if baseEvent.StreamID() != "5678" {
		t.Errorf("Expected StreamID %s, got %s", "5678", baseEvent.StreamID())
	}

	if baseEvent.StreamVersion() != 1 {
		t.Errorf("Expected StreamVersion 1, got %d", baseEvent.StreamVersion())
	}
}
