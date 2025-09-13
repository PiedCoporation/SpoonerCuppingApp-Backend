package abstractions

import "backend/internal/domains/entities"

type IEventUserRepository interface {
	GenericRepository[entities.EventUser]
}