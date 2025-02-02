package projection

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/messaging"
)

// Projector is a specific event listener that projects events into a read model
type Projector interface {
	Project(ctx context.Context, event domain.Event) error
	Subscribe(ctx context.Context, bus messaging.EventBus)
}
