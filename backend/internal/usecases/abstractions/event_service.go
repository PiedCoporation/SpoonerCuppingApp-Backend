package abstractions

import (
	"backend/internal/contracts/common"
	"backend/internal/contracts/event"
	"context"

	"github.com/google/uuid"
)

type IEventService interface {
	Create(ctx context.Context, req event.CreateEventReq) error
	GetAll(ctx context.Context, pageSize int, pageNumber int, searchTerm string) (*common.PageResult[event.Event], error)
	// GetByID(ctx context.Context, id uuid.UUID) (*entities.Event, error)
	
	// Update(ctx context.Context, id uuid.UUID, event *entities.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
}