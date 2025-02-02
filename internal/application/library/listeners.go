package library

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/logger"
	"github.com/mikaelchan/hamster/pkg/messaging"
)

func WhenLibraryCreated(cb messaging.CommandBus) domain.EventListener {
	return func(ctx context.Context, event domain.Event) error {
		logger.Infof("Library created: %s", event.ID().String())
		// TODO: scan the path, and create media items
		return nil
	}
}
