package abstractions

import "backend/internal/domains/entities"

type IEventAddressRepository interface {
	GenericRepository[entities.EventAddress]
}