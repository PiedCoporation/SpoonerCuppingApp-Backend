package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"

	"gorm.io/gorm"
)

type sampleTastingPgRepo struct {
	*genericRepository[entities.UserSampleTasting]
}

func NewSampleTastingRepo(db *gorm.DB) abstractions.ISampleTastingRepository {
	return &sampleTastingPgRepo{
		genericRepository: NewGenericRepository[entities.UserSampleTasting](db),
	}
}
