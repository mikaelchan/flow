package library

import (
	"context"

	"github.com/mikaelchan/hamster/internal/domain/shared"
	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/messaging"
	"github.com/mikaelchan/hamster/pkg/repository"
)

func OnCreateLibrary(idProvider domain.IDProvider, readModel ReadModel, repository repository.Repository) domain.CommandHandler {
	return func(ctx context.Context, cmd domain.Command) error {
		createLibraryCmd, ok := cmd.(*CreateLibrary)
		if !ok {
			return domain.ErrInvalidCommand
		}

		if err := createLibraryCmd.Validate(); err != nil {
			return err
		}

		// check if the library already exists
		exist, err := readModel.NameOrPathExists(ctx, createLibraryCmd.Name, createLibraryCmd.Location.Path)
		if err != nil {
			return err
		}
		if exist {
			return ErrLibraryAlreadyExists
		}

		// check if the location exists
		isWritable, err := createLibraryCmd.Location.IsWritable()
		if err != nil {
			return err
		}
		if !isWritable {
			return shared.ErrStorageLocationNotWritable
		}
		// check if the location has enough space
		freeSpace, err := createLibraryCmd.Location.FreeSpace()
		if err != nil {
			return err
		}
		if createLibraryCmd.Location.Capacity > freeSpace {
			return shared.ErrStorageLocationNotEnoughSpace
		}

		libID, err := idProvider.FetchID()
		if err != nil {
			return err
		}
		eventID, err := idProvider.FetchID()
		if err != nil {
			return err
		}
		library := &Library{}
		library.Create(libID, eventID, createLibraryCmd.Name, createLibraryCmd.MediaType, createLibraryCmd.Location, createLibraryCmd.QualityPreference, createLibraryCmd.NamingTemplate)

		return repository.Save(ctx, library)
	}
}

func Register(ctx context.Context, bus messaging.CommandBus, idProvider domain.IDProvider, readModel ReadModel, repository repository.Repository) {
	bus.Register(ctx, CreateLibraryContract, OnCreateLibrary(idProvider, readModel, repository))
}
