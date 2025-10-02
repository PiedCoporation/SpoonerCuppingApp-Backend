package abstractions

import "backend/internal/domains/entities"

type IEventSampleRepository interface {
	IGenericRepository[entities.EventSample]
}
