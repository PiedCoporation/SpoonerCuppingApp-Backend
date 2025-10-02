package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"

	"gorm.io/gorm"
)

type eventPgRepo struct {
	*genericRepository[entities.Event]
}

func NewEventRepo(db *gorm.DB) abstractions.IEventRepository {
	return &eventPgRepo{
		genericRepository: NewGenericRepository[entities.Event](db),
	}
}