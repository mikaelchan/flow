package postgresql

type EventData struct {
	StreamIDField string `gorm:"column:stream_id;not null;index:idx_events_stream_id"`
	EventIDField  string `gorm:"column:event_id;not null;primaryKey"`
	TypeField     string `gorm:"column:event_type;not null;index:idx_events_event_type"`
	VersionField  uint64 `gorm:"column:version;not null"`
	PayloadField  []byte `gorm:"column:payload;not null"`
}

func (EventData) TableName() string {
	return "events"
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
