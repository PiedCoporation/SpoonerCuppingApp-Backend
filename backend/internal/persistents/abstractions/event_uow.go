package abstractions

import "context"

type EventUOW interface {
	Begin(ctx context.Context) (EventRepoProvider, error)
	Commit() error
	Rollback() error
}

type EventRepoProvider interface {
	EventRepository() IEventRepository
	SampleRepository() ISampleRepository
	EventAddressRepository() IEventAddressRepository
	EventSampleRepository() IEventSampleRepository
}