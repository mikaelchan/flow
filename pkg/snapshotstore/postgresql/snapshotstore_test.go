package postgresql_test

import (
	"context"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mikaelchan/hamster/internal/domain/library"
	"github.com/mikaelchan/hamster/internal/domain/shared"
	_ "github.com/mikaelchan/hamster/internal/infrastructure/serializer"
	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/serializer"
	"github.com/mikaelchan/hamster/pkg/snapshotstore"
	"github.com/mikaelchan/hamster/pkg/snapshotstore/postgresql"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	dsn := "host=nas.myhome user=user password=password dbname=hamster port=5433 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestSaveAndLoad(t *testing.T) {
	policy := snapshotstore.SnapshotPolicy{
		ShouldSnapshot: func(ctx context.Context, aggregate domain.AggregateRoot) bool {
			return true
		},
	}
	ss := postgresql.NewPostgresSnapshotStore(db, serializer.GetFactory(), policy)

	t.Run("save", func(t *testing.T) {
		lib := library.Library{
			Name:      "test",
			MediaType: shared.Movie,
			Location: shared.StorageLocation{
				Path: "test",
			},
			QualityPreference: shared.QualityPreference{},
			NamingTemplate:    "test",
			Status:            library.Active,
		}
		lib.SetID("test")
		if err := ss.Save(context.Background(), &lib); err != nil {
			t.Fatalf("Failed to save snapshot: %v", err)
		}
	})

	t.Run("load", func(t *testing.T) {
		var loadedLibrary library.Library
		if err := ss.Load(context.Background(), "test", &loadedLibrary); err != nil {
			t.Fatalf("Failed to load snapshot: %v", err)
		}
		if loadedLibrary.ID() != "test" {
			t.Errorf("Expected ID %s, got %s", "test", loadedLibrary.ID())
		}
	})

	t.Run("delete", func(t *testing.T) {
		snapshot := postgresql.Snapshot{StreamIDField: "test"}
		//delete the snapshot
		if err := db.Delete(&snapshot).Error; err != nil {
			t.Fatalf("Failed to delete snapshot: %v", err)
		}
	})
}
