package abstractions

import "backend/internal/domains/entities"

type IEventUserRepository interface {
	IGenericRepository[entities.EventUser]
}
