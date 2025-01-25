package domain

type AggregateRoot interface {
	HasType
	ID() ID
	Version() Version
	UncommittedVersion() Version
	UncommittedEvents() []Event
	ClearUncommittedEvents()
	Append(event Event) error
	Apply(event Event) error
}

func Track(root AggregateRoot, event Event) error {
	if err := root.Apply(event); err != nil {
		return err
	}
	return root.Append(event)
}

type BaseAggregateRoot struct {
	id      ID
	version Version
	events  []Event
}

func (a *BaseAggregateRoot) ID() ID {
	return a.id
}

func (a *BaseAggregateRoot) Version() Version {
	return a.version
}

func (a *BaseAggregateRoot) UncommittedEvents() []Event {
	events := make([]Event, len(a.events))
	copy(events, a.events)
	return events
}

func (a *BaseAggregateRoot) ClearUncommittedEvents() {
	a.SetVersion(a.UncommittedVersion())
	a.events = []Event{}
}

func (a *BaseAggregateRoot) Apply(event Event) error {
	panic("not implemented")
}

func (a *BaseAggregateRoot) UncommittedVersion() Version {
	if len(a.events) == 0 {
		return a.version
	}
	return a.events[len(a.events)-1].StreamVersion()
}

func (a *BaseAggregateRoot) Append(event Event) error {
	a.events = append(a.events, event)
	return nil
}

func (a *BaseAggregateRoot) SetID(id ID) error {
	if a.id != EmptyID {
		return ErrIDAlreadySet
	}
	a.id = id
	return nil
}

func (a *BaseAggregateRoot) SetVersion(version Version) {
	a.version = version
}
