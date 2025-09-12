package abstractions

import "backend/internal/domains/entities"

type ISampleRepository interface {
	GenericRepository[entities.UserSample]
}