package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"

	"gorm.io/gorm"
)

type eventSamplePgRepo struct {
	*genericRepository[entities.EventSample]
}

func NewEventSampleRepo(db *gorm.DB) abstractions.IEventSampleRepository {
	return &eventSamplePgRepo{
		genericRepository: NewGenericRepository[entities.EventSample](db),
	}
}