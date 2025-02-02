package library

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mikaelchan/hamster/internal/domain/library"
	"github.com/mikaelchan/hamster/internal/infrastructure/idprovider/uuid"
	postgresql_readmodel "github.com/mikaelchan/hamster/internal/infrastructure/persistence/postgresql"
	_ "github.com/mikaelchan/hamster/internal/infrastructure/serializer"
	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/eventstore/postgresql"
	"github.com/mikaelchan/hamster/pkg/messaging"
	redismessaging "github.com/mikaelchan/hamster/pkg/messaging/redis"
	"github.com/mikaelchan/hamster/pkg/projection"
	"github.com/mikaelchan/hamster/pkg/repository"
	"github.com/mikaelchan/hamster/pkg/serializer"
	"github.com/mikaelchan/hamster/pkg/snapshotstore"
	postgresql_snapshotstore "github.com/mikaelchan/hamster/pkg/snapshotstore/postgresql"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service struct {
	db          *sql.DB
	redisClient *redis.Client

	readModel  library.ReadModel
	repository repository.Repository
	idProvider domain.IDProvider
	projector  projection.Projector
	commandBus messaging.CommandBus
	eventBus   messaging.EventBus
}

func NewService(ctx context.Context, config Config) *Service {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
	})
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Postgres.Host, config.Postgres.Port, config.Postgres.Username, config.Postgres.Password, config.Postgres.Database)
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db, err := gormDB.DB()
	if err != nil {
		panic(err)
	}
	factory := serializer.GetFactory()
	eb := redismessaging.NewEventBus(ctx,
		redismessaging.Config{
			Client:        redisClient,
			HandleTimeout: 5 * time.Second,
		},
		factory,
	)
	cb := redismessaging.NewCommandBus(ctx,
		redismessaging.Config{
			Client:        redisClient,
			HandleTimeout: 5 * time.Second,
		},
		factory,
	)
	es := postgresql.NewEventStore(gormDB, factory)
	ss := postgresql_snapshotstore.NewPostgresSnapshotStore(gormDB, factory, snapshotstore.SnapshotPolicy{
		ShouldSnapshot: func(ctx context.Context, root domain.AggregateRoot) bool {
			return root.Version()%10 == 0
		},
	})
	idProvider := uuid.NewIDProvider()
	repo := repository.NewSnapshotRepository(es, ss, eb)
	readModel := postgresql_readmodel.NewLibraryReadModel(gormDB)
	projector := NewProjector(readModel)
	projector.Subscribe(ctx, eb)
	return &Service{db: db, redisClient: redisClient, idProvider: idProvider, readModel: readModel, repository: repo, projector: projector, commandBus: cb, eventBus: eb}
}

func (s *Service) Close() error {
	var firstErr error
	if err := s.commandBus.Close(); err != nil {
		firstErr = err
	}
	if err := s.eventBus.Close(); err != nil && firstErr == nil {
		firstErr = err
	}
	if err := s.redisClient.Close(); err != nil && firstErr == nil {
		firstErr = err
	}
	if err := s.db.Close(); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}
