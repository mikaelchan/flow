package eventstore

import "errors"

var (
	ErrNoEvents       = errors.New("no events")
	ErrStreamNotFound = errors.New("stream not found")
)
