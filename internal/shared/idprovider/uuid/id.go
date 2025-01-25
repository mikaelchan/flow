package uuid

import (
	"github.com/google/uuid"

	"github.com/mikaelchan/hamster/pkg/domain"
)

type idProvider struct{}

func NewIDProvider() domain.IDProvider {
	uuid.EnableRandPool()
	return &idProvider{}
}

func (p *idProvider) FetchID() (domain.ID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return domain.EmptyID, err
	}
	return domain.ID(id.String()), nil
}
