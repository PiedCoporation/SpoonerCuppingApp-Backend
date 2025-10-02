package abstractions

import "backend/internal/domains/entities"

type IPostRepository interface {
	IGenericRepository[entities.Post]
}
