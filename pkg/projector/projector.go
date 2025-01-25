package projector

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
)

type ReadModel interface {
	Apply(ctx context.Context, events domain.Event) error
	Version() domain.Version
	Flush(ctx context.Context) error
}

// Projector defines the methods for projecting events to read models
// type Projector interface {
// 	Run(ctx context.Context) error
// 	Close(ctx context.Context) error
// 	Project(ctx context.Context, events domain.Event) error
// }

// type projector struct {
// 	running       bool
// 	trigger       chan domain.Event
// 	stop          chan struct{}
// 	flushInterval time.Duration
// 	readModel     ReadModel
// }

// func NewProjector(readModel ReadModel, flushInterval time.Duration) Projector {
// 	return &projector{
// 		readModel:     readModel,
// 		trigger:       make(chan domain.Event),
// 		stop:          make(chan struct{}),
// 		flushInterval: flushInterval,
// 	}
// }

// func (p *projector) Run(ctx context.Context) error {
// 	p.running = true
// 	go func() {
// 		select {
// 		case <-p.stop:
// 			p.running = false
// 		}
// 	}()
// 	return nil
// }

// func (p *projector) Close(ctx context.Context) error {
// 	return nil
// }
