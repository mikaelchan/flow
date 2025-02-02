package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/mikaelchan/hamster/internal/domain/library"
	"github.com/mikaelchan/hamster/internal/domain/shared"
	"github.com/mikaelchan/hamster/pkg/domain"
	"gorm.io/gorm"
)

type Library struct {
	ID                string            `gorm:"primaryKey"`
	Name              string            `gorm:"not null;unique"`
	MediaType         uint8             `gorm:"not null"`
	Location          string            `gorm:"not null"`
	QualityPreference map[string]string `gorm:"not null;serializer:json"`
	NamingTemplate    string            `gorm:"not null"`
	Status            uint8             `gorm:"not null"`
	CreatedAt         time.Time         `gorm:"not null"`
	UpdatedAt         time.Time         `gorm:"not null"`
}

func (Library) TableName() string {
	return "libraries"
}

type LibraryReadModel struct {
	db *gorm.DB
}

func NewLibraryReadModel(db *gorm.DB) library.ReadModel {
	return &LibraryReadModel{db: db}
}

// Delete implements library.ReadModel.
func (m *LibraryReadModel) Delete(ctx context.Context, libraryID domain.ID) error {
	return m.db.Delete(&Library{}, "id = ?", libraryID.String()).Error
}

// NameOrPathExists implements library.ReadModel.
func (m *LibraryReadModel) NameOrPathExists(ctx context.Context, name string, path string) (bool, error) {
	var lib Library
	if err := m.db.Where("name = ?", name).Or("location = ?", path).First(&lib).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// UpdateNamingTemplate implements library.ReadModel.
func (m *LibraryReadModel) UpdateNamingTemplate(ctx context.Context, libraryID domain.ID, namingTemplate shared.NamingTemplate, updatedAt time.Time) error {
	return m.db.Model(&Library{}).Where("id = ?", libraryID.String()).Update("naming_template", namingTemplate.String()).Update("updated_at", updatedAt).Error
}

// UpdateQualityPreference implements library.ReadModel.
func (m *LibraryReadModel) UpdateQualityPreference(ctx context.Context, libraryID domain.ID, qualityPreference shared.QualityPreference, updatedAt time.Time) error {
	return m.db.Model(&Library{}).Where("id = ?", libraryID.String()).Update("quality_preference", qualityPreference).Update("updated_at", updatedAt).Error
}

// UpdateStatus implements library.ReadModel.
func (m *LibraryReadModel) UpdateStatus(ctx context.Context, libraryID domain.ID, status library.Status, updatedAt time.Time) error {
	return m.db.Model(&Library{}).Where("id = ?", libraryID.String()).Update("status", uint8(status)).Update("updated_at", updatedAt).Error
}

func (m *LibraryReadModel) Create(ctx context.Context, id domain.ID, name string, mediaType shared.MediaType, location shared.StorageLocation, qualityPreference shared.QualityPreference, namingTemplate shared.NamingTemplate, createdAt time.Time) error {
	lib := &Library{
		ID:                id.String(),
		Name:              name,
		MediaType:         uint8(mediaType),
		Location:          string(location.Path),
		QualityPreference: qualityPreference,
		NamingTemplate:    string(namingTemplate),
		Status:            uint8(library.Active),
		CreatedAt:         createdAt,
		UpdatedAt:         createdAt,
	}
	return m.db.Create(&lib).Error
}
