package projection

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/messaging"
)

type Projector interface {
	Project(ctx context.Context, event domain.Event) error
	Subscribe(ctx context.Context, bus messaging.EventBus)
}
