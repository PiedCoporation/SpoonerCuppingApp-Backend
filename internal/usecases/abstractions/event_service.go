package abstractions

import (
	"backend/internal/contracts/common"
	"backend/internal/contracts/event"
	"context"

	"github.com/google/uuid"
)

type IEventService interface {
    Create(ctx context.Context, req event.CreateEventReq) (*common.Result[event.Event])
	GetAll(ctx context.Context, pageSize int, pageNumber int, searchTerm string) (*common.Result[common.PageResult[event.Event]])
	GetByID(ctx context.Context, id uuid.UUID) (*common.Result[event.GetEventByIDResponse])
	Register(ctx context.Context, id uuid.UUID) (*common.Result[string])
	StartEvent(ctx context.Context, id uuid.UUID) (*common.Result[string])
	// Update(ctx context.Context, id uuid.UUID, event *entities.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetEventByUserID(ctx context.Context, pageSize int, pageNumber int) (*common.Result[common.PageResult[event.Event]])
}