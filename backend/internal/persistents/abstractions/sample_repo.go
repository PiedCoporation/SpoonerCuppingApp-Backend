package abstractions

import "backend/internal/domains/entities"

type ISampleRepository interface {
	IGenericRepository[entities.UserSample]
}
