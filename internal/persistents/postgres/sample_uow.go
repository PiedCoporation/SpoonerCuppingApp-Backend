package postgres

import (
	"backend/internal/persistents/abstractions"
	"context"

	"gorm.io/gorm"
)

type sampleUow struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewSampleUow(db *gorm.DB) abstractions.SampleUOW {
	return &sampleUow{db: db}
}

func (u *sampleUow) Begin(ctx context.Context) (abstractions.SampleRepoProvider, error) {
	tx := u.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	u.tx = tx
	return &sampleRepoProvider{tx: tx}, nil
}

func (u *sampleUow) Commit() error {
	if u.tx == nil {
		return nil
	}
	return u.tx.Commit().Error
}

func (u *sampleUow) Rollback() error {
	if u.tx == nil {
		return nil
	}
	return u.tx.Rollback().Error
}

func (u *sampleUow) GetDB() *gorm.DB {
	return u.db
}

type sampleRepoProvider struct {
	tx *gorm.DB
	sampleRepo abstractions.ISampleRepository
	sampleTastingRepo abstractions.ISampleTastingRepository
}

func (r *sampleRepoProvider) SampleRepository() abstractions.ISampleRepository {
	if r.sampleRepo == nil {
		r.sampleRepo = NewSampleRepo(r.tx)
	}
	return r.sampleRepo
}

func (r *sampleRepoProvider) SampleTastingRepository() abstractions.ISampleTastingRepository {
	if r.sampleTastingRepo == nil {
		r.sampleTastingRepo = NewSampleTastingRepo(r.tx)
	}
	return r.sampleTastingRepo
}