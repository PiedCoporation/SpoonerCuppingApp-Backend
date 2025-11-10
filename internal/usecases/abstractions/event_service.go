package abstractions

import (
	"backend/internal/constants/enums/eventparticipant"
	"backend/internal/contracts/common"
	eventContractRequest "backend/internal/contracts/event/request"
	eventContractResponse "backend/internal/contracts/event/response"
	"context"

	"github.com/google/uuid"
)

type IEventService interface {
    Create(ctx context.Context, req eventContractRequest.CreateEventReq) (*common.Result[eventContractResponse.Event])
	GetAll(ctx context.Context, pageSize int, pageNumber int, searchTerm string) (*common.Result[common.PageResult[eventContractResponse.Event]])
	GetByID(ctx context.Context, id uuid.UUID) (*common.Result[eventContractResponse.GetEventByIDResponse])
	Register(ctx context.Context, id uuid.UUID) (*common.Result[string])
	ResponseEvent(ctx context.Context, userID uuid.UUID, eventUserID uuid.UUID, isAccept bool) (*common.Result[string])
	GetEventParticipant(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, pageSize int, pageNumber int, searchTerm string, typeParticipant eventparticipant.TypeParticipantEnum) (*common.Result[common.PageResult[eventContractResponse.GetAllUserEventResponse]])
	StartEvent(ctx context.Context, id uuid.UUID) (*common.Result[string])
	// Update(ctx context.Context, id uuid.UUID, event *entities.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetEventByUserID(ctx context.Context, pageSize int, pageNumber int) (*common.Result[common.PageResult[eventContractResponse.Event]])
}