package library

type Status uint8

const (
	Active Status = iota
	Inactive
	Archived
)
