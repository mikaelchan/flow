package repository

import (
	"context"

	"github.com/mikaelchan/hamster/pkg/domain"
)

type Repository interface {
	Save(ctx context.Context, root domain.AggregateRoot) error
	Load(ctx context.Context, id domain.ID, root domain.AggregateRoot) error
}
