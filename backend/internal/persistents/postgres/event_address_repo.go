package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"

	"gorm.io/gorm"
)

type eventAddressPgRepo struct {
	*genericRepository[entities.EventAddress]
}

func NewEventAddressRepo(db *gorm.DB) abstractions.IEventAddressRepository {
	return &eventAddressPgRepo{
		genericRepository: NewGenericRepository[entities.EventAddress](db),
	}
}