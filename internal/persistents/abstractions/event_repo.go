package abstractions

import "backend/internal/domains/entities"

type IEventRepository interface {
	IGenericRepository[entities.Event]
}
