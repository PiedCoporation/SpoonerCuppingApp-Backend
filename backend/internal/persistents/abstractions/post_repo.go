package abstractions

import "backend/internal/domains/entities"

type IPostRepository interface {
	GenericRepository[entities.Post]
}
