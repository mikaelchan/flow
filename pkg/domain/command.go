package domain

import "context"

// CommandHandler is a function that handles a command
type CommandHandler func(ctx context.Context, cmd Command) error

// Command represents a command in the system
type Command interface {
	HasType
	// Validate validates the command
	Validate() error
}
