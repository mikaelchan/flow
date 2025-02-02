package postgresql

type Snapshot struct {
	StreamIDField string `gorm:"column:stream_id;primaryKey"`
	TypeField     string `gorm:"column:type"`
	VersionField  uint64 `gorm:"column:version"`
	StateField    []byte `gorm:"column:state"`
}

func (s *Snapshot) StreamID() string {
	return s.StreamIDField
}

func (s *Snapshot) Type() string {
	return s.TypeField
}

func (s *Snapshot) Version() uint64 {
	return s.VersionField
}

func (s *Snapshot) State() []byte {
	return s.StateField
}

func (s *Snapshot) TableName() string {
	return "snapshots"
}
