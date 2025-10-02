package abstractions

import (
	"context"

	"gorm.io/gorm"
)

type EventUOW interface {
	Begin(ctx context.Context) (EventRepoProvider, error)
	Commit() error
	Rollback() error
	GetDB() *gorm.DB
}

type EventRepoProvider interface {
	EventRepository() IEventRepository
	SampleRepository() ISampleRepository
	EventAddressRepository() IEventAddressRepository
	EventSampleRepository() IEventSampleRepository
	EventUserRepository() IEventUserRepository
}