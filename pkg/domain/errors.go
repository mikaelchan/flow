package domain

import "errors"

var (
	ErrHandlerNotFound = errors.New("no handler registered for command")
	ErrInvalidCommand  = errors.New("invalid command")
	ErrCommandTimeout  = errors.New("command handling timed out")
	ErrIDAlreadySet    = errors.New("aggregate id already set")
)
