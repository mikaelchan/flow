package library

import (
	"time"

	"github.com/mikaelchan/hamster/internal/domain/shared"
	"github.com/mikaelchan/hamster/pkg/domain"
)

const (
	CreatedEventTopic domain.Type = "library.created"
)

type Created struct {
	domain.BaseEvent
	LibraryID         domain.ID                `json:"library_id"`
	Name              string                   `json:"name"`
	MediaType         shared.MediaType         `json:"media_type"`
	Location          shared.StorageLocation   `json:"location"`
	QualityPreference shared.QualityPreference `json:"quality_preference"`
	NamingTemplate    shared.NamingTemplate    `json:"naming_template"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

func (e *Created) Type() domain.Type {
	return CreatedEventTopic
}
