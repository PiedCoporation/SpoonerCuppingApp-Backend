package abstractions

import "backend/internal/domains/entities"

type IEventRepository interface {
	GenericRepository[entities.Event]
}