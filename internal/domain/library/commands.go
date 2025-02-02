package library

import (
	"github.com/mikaelchan/hamster/internal/domain/shared"
	"github.com/mikaelchan/hamster/pkg/domain"
)

const (
	CreateLibraryCommandType domain.Type = "create-library"
)

type CreateLibrary struct {
	Name              string                   `json:"name"`
	MediaType         shared.MediaType         `json:"media_type"`
	Location          shared.StorageLocation   `json:"location"`
	QualityPreference shared.QualityPreference `json:"quality_preference"`
	NamingTemplate    shared.NamingTemplate    `json:"naming_template"`
}

func (cmd *CreateLibrary) Type() domain.Type {
	return CreateLibraryCommandType
}

func (cmd *CreateLibrary) Validate() error {
	// TODO: add more validation
	if cmd.Name == "" {
		return ErrEmptyName
	}
	if !cmd.MediaType.IsValid() {
		return shared.ErrInvalidMediaType
	}
	if cmd.Location.Path == "" {
		return shared.ErrEmptyPath
	}
	if _, err := shared.NewNamingTemplateParser(cmd.NamingTemplate); err != nil {
		return shared.ErrInvalidNamingTemplate
	}
	return nil
}
