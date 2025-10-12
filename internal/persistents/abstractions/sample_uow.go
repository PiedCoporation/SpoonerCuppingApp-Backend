package abstractions

import (
	"context"

	"gorm.io/gorm"
)

type SampleUOW interface {
	Begin(ctx context.Context) (SampleRepoProvider, error)
	Commit() error
	Rollback() error
	GetDB() *gorm.DB
}

type SampleRepoProvider interface {
	SampleRepository() ISampleRepository
	SampleTastingRepository() ISampleTastingRepository
}