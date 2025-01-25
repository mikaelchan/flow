package domain

type Version uint64

type ID string

const EmptyID = ID("")

func (id ID) String() string {
	return string(id)
}

type Type string

func (t Type) String() string {
	return string(t)
}

type HasType interface {
	Type() Type
}

type IDProvider interface {
	FetchID() (ID, error)
}
