package postgresql

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/eventstore"
	"github.com/mikaelchan/hamster/pkg/serializer"
	"gorm.io/gorm"
)

type EventStore struct {
	db      *gorm.DB
	factory *serializer.Factory
}

func NewEventStore(db *gorm.DB, factory *serializer.Factory) eventstore.EventStore {
	return &EventStore{
		db:      db,
		factory: factory,
	}
}

func (s *EventStore) Append(ctx context.Context, streamID domain.ID, events ...domain.Event) error {
	if len(events) == 0 {
		return eventstore.ErrNoEvents
	}

	var eventDatas []EventData
	for _, event := range events {
		payload, err := s.factory.Serialize(event)
		if err != nil {
			return err
		}
		eventDatas = append(eventDatas, EventData{
			StreamIDField: streamID.String(),
			EventIDField:  event.ID().String(),
			TypeField:     string(event.Type()),
			VersionField:  uint64(event.StreamVersion()),
			PayloadField:  payload,
		})
	}

	return s.db.Create(&eventDatas).Error
}

func (s *EventStore) Load(ctx context.Context, streamID domain.ID, from domain.Version) (eventstore.EventIterator, error) {
	var eventDatas []EventData
	if err := s.db.Where("stream_id = ? AND version >= ?", streamID, from).Find(&eventDatas).Error; err != nil {
		return nil, err
	}

	return &eventIterator{eventDatas: eventDatas, factory: s.factory}, nil
}
