package mongo

import (
	"github.com/mikaelchan/hamster/pkg/eventstore"
)

var _ eventstore.EventData = &EventData{}

type EventData struct {
	StreamIDField string `bson:"stream_id"`
	EventIDField  string `bson:"event_id"`
	TypeField     string `bson:"event_type"`
	VersionField  uint64 `bson:"version"`
	PayloadField  []byte `bson:"payload"`
}

func (m *EventData) StreamID() string {
	return m.StreamIDField
}

func (m *EventData) Type() string {
	return m.TypeField
}

func (m *EventData) Payload() []byte {
	return m.PayloadField
}
