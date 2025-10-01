package abstractions

import "backend/internal/domains/entities"

type IEventAddressRepository interface {
	IGenericRepository[entities.EventAddress]
}
