package postgres

import (
	"backend/internal/persistents/abstractions"
	"context"

	"gorm.io/gorm"
)

type eventUow struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewEventUow(db *gorm.DB) abstractions.EventUOW {
	return &eventUow{db: db}
}


func (u *eventUow) Begin(ctx context.Context) (abstractions.EventRepoProvider, error) {
	tx := u.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	u.tx = tx
	return &eventRepoProvider{tx: tx}, nil
}

func (u *eventUow) Commit() error {
	if u.tx == nil {
		return nil
	}
	return u.tx.Commit().Error
}

func (u *eventUow) Rollback() error {
	if u.tx == nil {
		return nil
	}
	return u.tx.Rollback().Error
}

type eventRepoProvider struct {
	tx *gorm.DB
	eventRepo abstractions.IEventRepository
	sampleRepo abstractions.ISampleRepository
	eventAddressRepo abstractions.IEventAddressRepository
}

func (r *eventRepoProvider) EventRepository() abstractions.IEventRepository {
	if r.eventRepo == nil {
		r.eventRepo = NewEventRepo(r.tx)
	}
	return r.eventRepo	
}

func (r *eventRepoProvider) SampleRepository() abstractions.ISampleRepository {
	if r.sampleRepo == nil {
		r.sampleRepo = NewSampleRepo(r.tx)
	}
	return r.sampleRepo
}

func (r *eventRepoProvider) EventAddressRepository() abstractions.IEventAddressRepository {
	if r.eventAddressRepo == nil {
		r.eventAddressRepo = NewEventAddressRepo(r.tx)
	}
	return r.eventAddressRepo
}