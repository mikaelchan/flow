package domain

import (
	"context"
)

// EventListener is a function that listens to an event
type EventListener func(ctx context.Context, event Event) error

// Event represents a domain event
type Event interface {
	HasType

	// ID returns the ID of the event
	ID() ID

	// StreamID returns the ID of the stream that produced the event
	StreamID() ID

	// StreamVersion returns the version of the stream that produced the event
	StreamVersion() Version
}

type BaseEvent struct {
	EventID            ID      `json:"id"`
	EventStreamID      ID      `json:"stream_id"`
	EventStreamVersion Version `json:"stream_version"`
}

func (e *BaseEvent) ID() ID {
	return e.EventID
}

func (e *BaseEvent) StreamID() ID {
	return e.EventStreamID
}

func (e *BaseEvent) StreamVersion() Version {
	return e.EventStreamVersion
}

func NewBaseEvent(id ID, root AggregateRoot) BaseEvent {
	return BaseEvent{
		EventID:            id,
		EventStreamID:      root.ID(),
		EventStreamVersion: root.UncommittedVersion() + 1,
	}
}
