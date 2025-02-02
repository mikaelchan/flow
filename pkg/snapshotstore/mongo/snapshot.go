package mongo

// Snapshot is a concrete implementation of the Snapshot interface
type Snapshot struct {
	StreamIDField string `bson:"stream_id"`
	TypeField     string `bson:"type"`
	VersionField  uint64 `bson:"version"`
	StateField    []byte `bson:"state"`
}

// StreamID returns the stream ID of the snapshot
func (s *Snapshot) StreamID() string {
	return s.StreamIDField
}

// Type returns the type of the snapshot
func (s *Snapshot) Type() string {
	return s.TypeField
}

// Version returns the version of the snapshot
func (s *Snapshot) Version() uint64 {
	return s.VersionField
}

// State returns the state of the snapshot
func (s *Snapshot) State() []byte {
	return s.StateField
}
