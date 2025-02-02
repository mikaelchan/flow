package library

import (
	"context"
	"time"

	"github.com/mikaelchan/hamster/internal/domain/shared"
	"github.com/mikaelchan/hamster/pkg/domain"
)

type ReadModel interface {
	// read methods
	NameOrPathExists(ctx context.Context, name string, path string) (bool, error)
	// write methods
	Create(ctx context.Context, libraryID domain.ID, name string, mediaType shared.MediaType, location shared.StorageLocation, qualityPreference shared.QualityPreference, namingTemplate shared.NamingTemplate, createdAt time.Time) error
	UpdateQualityPreference(ctx context.Context, libraryID domain.ID, qualityPreference shared.QualityPreference, updatedAt time.Time) error
	UpdateNamingTemplate(ctx context.Context, libraryID domain.ID, namingTemplate shared.NamingTemplate, updatedAt time.Time) error
	UpdateStatus(ctx context.Context, libraryID domain.ID, status Status, updatedAt time.Time) error
	Delete(ctx context.Context, libraryID domain.ID) error
}
