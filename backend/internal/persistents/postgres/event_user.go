package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"

	"gorm.io/gorm"
)

type eventUserRepo struct {
	*genericRepository[entities.EventUser]
}

func NewEventUserRepo(db *gorm.DB) abstractions.IEventUserRepository {
	return &eventUserRepo{
		genericRepository: NewGenericRepository[entities.EventUser](db),
	}
}