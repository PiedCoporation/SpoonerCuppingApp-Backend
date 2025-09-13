package abstractions

import "backend/internal/domains/entities"

type IEventSampleRepository interface {
	GenericRepository[entities.EventSample]
}