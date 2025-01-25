package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/eventstore"
	"github.com/mikaelchan/hamster/pkg/serializer"
)

var _ eventstore.EventIterator = &eventIterator{}

type eventIterator struct {
	cursor  *mongo.Cursor
	factory *serializer.Factory
}

func (m *eventIterator) HasNext(ctx context.Context) bool {
	return m.cursor.Next(ctx)
}

func (m *eventIterator) Next(ctx context.Context) (domain.Event, error) {
	var eventData EventData
	if err := m.cursor.Decode(&eventData); err != nil {
		return nil, err
	}

	event, err := m.factory.DeserializeNew(domain.Type(eventData.Type()), eventData.Payload())
	if err != nil {
		return nil, err
	}
	return event.(domain.Event), nil
}

func (m *eventIterator) Close(ctx context.Context) error {
	return m.cursor.Close(ctx)
}
