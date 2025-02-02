package postgresql

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/eventstore"
	"github.com/mikaelchan/hamster/pkg/serializer"
)

var _ eventstore.EventIterator = &eventIterator{}

type eventIterator struct {
	eventDatas []EventData
	factory    *serializer.Factory
	index      int
}

func (m *eventIterator) HasNext(ctx context.Context) bool {
	return m.index < len(m.eventDatas)
}

func (m *eventIterator) Next(ctx context.Context) (domain.Event, error) {
	eventData := m.eventDatas[m.index]
	m.index++

	event, err := m.factory.DeserializeNew(domain.Type(eventData.Type()), eventData.Payload())
	if err != nil {
		return nil, err
	}
	return event.(domain.Event), nil
}

func (m *eventIterator) Close(ctx context.Context) error {
	return nil
}
