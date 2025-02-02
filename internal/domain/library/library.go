package library

import (
	"time"

	"github.com/mikaelchan/hamster/internal/domain/shared"
	"github.com/mikaelchan/hamster/pkg/domain"
)

type Library struct {
	domain.BaseAggregateRoot
	Name              string
	MediaType         shared.MediaType
	Location          shared.StorageLocation
	QualityPreference shared.QualityPreference
	NamingTemplate    shared.NamingTemplate
	Status            Status

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (lib *Library) Type() domain.Type {
	return "library"
}

func (lib *Library) Create(id domain.ID, eventID domain.ID, name string, mediaType shared.MediaType, location shared.StorageLocation, qualityPreference shared.QualityPreference, namingTemplate shared.NamingTemplate) error {
	event := &Created{
		BaseEvent:         domain.NewBaseEvent(eventID, lib),
		LibraryID:         id,
		Name:              name,
		MediaType:         mediaType,
		Location:          location,
		QualityPreference: qualityPreference,
		NamingTemplate:    namingTemplate,
	}

	return domain.Track(lib, event)
}

func (lib *Library) Apply(event domain.Event) error {
	switch event := event.(type) {
	case *Created:
		lib.SetID(event.LibraryID)
		lib.Name = event.Name
		lib.MediaType = event.MediaType
		lib.Location = event.Location
		lib.QualityPreference = event.QualityPreference
		lib.NamingTemplate = event.NamingTemplate
		lib.Status = Active
		lib.CreatedAt = event.CreatedAt
		lib.UpdatedAt = event.UpdatedAt
	}
	return nil
}
