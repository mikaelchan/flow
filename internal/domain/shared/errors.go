package shared

import "errors"

var (
	ErrInvalidMediaType              = errors.New("invalid media type")
	ErrEmptyPath                     = errors.New("path cannot be empty")
	ErrInvalidCapacity               = errors.New("invalid capacity")
	ErrInvalidNamingTemplate         = errors.New("invalid naming template")
	ErrStorageLocationNotFound       = errors.New("storage location not found")
	ErrStorageLocationNotWritable    = errors.New("storage location not writable")
	ErrStorageLocationNotEnoughSpace = errors.New("storage location not enough space")
)
