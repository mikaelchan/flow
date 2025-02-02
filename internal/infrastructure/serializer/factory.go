package serializer

import (
	"encoding/json"
	"time"

	"github.com/mikaelchan/hamster/internal/domain/library"
	"github.com/mikaelchan/hamster/internal/domain/shared"
	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/serializer"
	jsonserializer "github.com/mikaelchan/hamster/pkg/serializer/json"
)

type Library struct {
	ID                string                   `json:"id"`
	Name              string                   `json:"name"`
	MediaType         shared.MediaType         `json:"media_type"`
	Location          shared.StorageLocation   `json:"location"`
	QualityPreference shared.QualityPreference `json:"quality_preference"`
	NamingTemplate    shared.NamingTemplate    `json:"naming_template"`
	Status            library.Status           `json:"status"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

func init() {
	factory := serializer.GetFactory()
	jsonserializer.Register(factory, &library.Created{}, &library.CreateLibrary{})

	factory.Register(&library.Library{}, func(val domain.HasType) ([]byte, error) {
		lib := val.(*library.Library)
		return json.Marshal(Library{
			ID:                lib.ID().String(),
			Name:              lib.Name,
			MediaType:         lib.MediaType,
			Location:          lib.Location,
			QualityPreference: lib.QualityPreference,
			NamingTemplate:    lib.NamingTemplate,
			Status:            lib.Status,
			CreatedAt:         lib.CreatedAt,
			UpdatedAt:         lib.UpdatedAt,
		})
	}, func(data []byte, val domain.HasType) error {
		lib := val.(*library.Library)
		var libData Library
		if err := json.Unmarshal(data, &libData); err != nil {
			return err
		}

		if err := lib.SetID(domain.ID(libData.ID)); err != nil {
			return err
		}
		lib.Name = libData.Name
		lib.MediaType = libData.MediaType
		lib.Location = libData.Location
		lib.QualityPreference = libData.QualityPreference
		lib.NamingTemplate = libData.NamingTemplate
		lib.Status = libData.Status
		lib.CreatedAt = libData.CreatedAt
		lib.UpdatedAt = libData.UpdatedAt
		return nil
	})
}
