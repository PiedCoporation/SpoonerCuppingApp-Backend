package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"

	"gorm.io/gorm"
)

type samplePgRepo struct {
	*genericRepository[entities.UserSample]
}

func NewSampleRepo(db *gorm.DB) abstractions.ISampleRepository {
	return &samplePgRepo{
		genericRepository: NewGenericRepository[entities.UserSample](db),
	}
}