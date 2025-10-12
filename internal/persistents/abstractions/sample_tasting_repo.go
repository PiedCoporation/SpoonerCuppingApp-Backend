package abstractions

import "backend/internal/domains/entities"

type ISampleTastingRepository interface {
	IGenericRepository[entities.UserSampleTasting]
}
