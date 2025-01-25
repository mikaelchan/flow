package eventstore

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
)

var _ EventIterator = &EmptyIterator{}

type EmptyIterator struct{}

func (e *EmptyIterator) HasNext(ctx context.Context) bool {
	return false
}

func (e *EmptyIterator) Next(ctx context.Context) (domain.Event, error) {
	return nil, ErrNoEvents
}

func (e *EmptyIterator) Close(ctx context.Context) error {
	return nil
}
