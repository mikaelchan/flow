package library

import (
	"context"

	"github.com/mikaelchan/hamster/internal/domain/library"
	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/messaging"
	"github.com/mikaelchan/hamster/pkg/projection"
)

type projector struct {
	readModel library.ReadModel
}

func NewProjector(readModel library.ReadModel) projection.Projector {
	return &projector{readModel: readModel}
}

func (p *projector) Project(ctx context.Context, event domain.Event) error {
	switch event := event.(type) {
	case *library.Created:
		return p.readModel.Create(ctx, event.LibraryID, event.Name, event.MediaType, event.Location, event.QualityPreference, event.NamingTemplate, event.CreatedAt)
		//TODO: handle other events
	}
	return nil
}

func (p *projector) Subscribe(ctx context.Context, bus messaging.EventBus) {
	listener := func(ctx context.Context, event domain.Event) error {
		return p.Project(ctx, event)
	}
	bus.Subscribe(ctx, library.CreatedEventTopic, listener)
	// bus.Subscribe(ctx, "library-updated", listener)
	// bus.Subscribe(ctx, "library-deleted", listener)
}
