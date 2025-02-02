package library

import "errors"

var (
	ErrEmptyName            = errors.New("name cannot be empty")
	ErrLibraryAlreadyExists = errors.New("library already exists")
)
