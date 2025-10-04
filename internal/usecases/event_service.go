package usecases

import (
	"backend/global"
	"backend/internal/constants/enums/eventregisterstatus"
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/common"
	"backend/internal/contracts/event"
	"backend/internal/domains/commons"
	"backend/internal/domains/entities"
	"backend/internal/mapper"
	persistentRepo "backend/internal/persistents/abstractions"
	"backend/internal/persistents/postgres"
	abstractions "backend/internal/usecases/abstractions"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type eventService struct {
	eventUOW persistentRepo.EventUOW
	eventRepo persistentRepo.IEventRepository
	eventAddressRepo persistentRepo.IEventAddressRepository
	eventSampleRepo persistentRepo.IEventSampleRepository
	eventUserRepo persistentRepo.IEventUserRepository
}

func NewEventService(eventUOW persistentRepo.EventUOW,
	 eventRepo persistentRepo.IEventRepository,
	 eventAddressRepo persistentRepo.IEventAddressRepository,
	 eventSampleRepo persistentRepo.IEventSampleRepository,
	 eventUserRepo persistentRepo.IEventUserRepository) abstractions.IEventService {
	return &eventService{
		eventUOW: eventUOW,
		eventRepo: eventRepo,
		eventAddressRepo: eventAddressRepo,
		eventSampleRepo: eventSampleRepo,
		eventUserRepo: eventUserRepo,
	}
}

func (s *eventService) Register(ctx context.Context, id uuid.UUID) error {
	userID, _ := ctx.Value("userID").(uuid.UUID)

	global.Logger.Info("userID", zap.Any("userID", userID))

	repoProvider, err := s.eventUOW.Begin(ctx)
	if err != nil {
		return err
	}

	eventRepo := repoProvider.EventRepository()
	eventUserRepo := repoProvider.EventUserRepository()

	eventEntity, err := eventRepo.GetByID(ctx, id)
	if err != nil {
		s.eventUOW.Rollback()
		return errorcode.ErrEventNotFound
	}

	if eventEntity.RegisterStatus == eventregisterstatus.RegisterStatusEnumPending {
		s.eventUOW.Rollback()
		return errorcode.ErrEventIsNotStartForRegister
	}

	if eventEntity.RegisterStatus == eventregisterstatus.RegisterStatusEnumFull {
		s.eventUOW.Rollback()
		return errorcode.ErrEventIsFull
	}

	if eventEntity.RegisterDate.After(time.Now()) {
		s.eventUOW.Rollback()
		return errorcode.ErrEventIsNotStartForRegister
	}

	if eventEntity.TotalCurrent >= eventEntity.Limit {
		s.eventUOW.Rollback()
		return errorcode.ErrEventIsFull
	}

	query := fmt.Sprintf("user_id = '%s' AND event_id = '%s'", userID.String(), id.String())

	// Check if user is already registered for this event
	existingEventUser, err := eventUserRepo.GetSingle(ctx, query)
	if err != nil {
		s.eventUOW.Rollback()
		return err
	}
	if existingEventUser != nil {
		s.eventUOW.Rollback()
		return errorcode.ErrUserAlreadyRegistered
	}

	if err := eventUserRepo.Create(ctx, &entities.EventUser{
		Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
		UserID: userID,
		EventID: id,
		IsAccepted: false,
		IsInvited: false,
	}); err != nil {
		s.eventUOW.Rollback()
		return err
	}

	eventEntity.TotalCurrent++
	if err := eventRepo.Update(ctx, eventEntity.ID, map[string]any{
		"total_current": eventEntity.TotalCurrent,
	}); err != nil {
		s.eventUOW.Rollback()
		return err
	}

	if err := s.eventUOW.Commit(); err != nil {
		s.eventUOW.Rollback()
		return err
	}

	return nil
}

func (s *eventService) Create(ctx context.Context, req event.CreateEventReq) (*common.Result[event.Event]) {
	// g, gCtx := errgroup.WithContext(ctx)
	userID, _ := ctx.Value("userID").(uuid.UUID)

	repoProvider, err := s.eventUOW.Begin(ctx)
    if err != nil {
        return common.Failure[event.Event](&common.Error{Code: "500", Message: "Failed to begin transaction"})
    }

	sampleRepo := repoProvider.SampleRepository()
	eventRepo := repoProvider.EventRepository()
	eventAddressRepo := repoProvider.EventAddressRepository()
	eventSampleRepo := repoProvider.EventSampleRepository()

    if len(req.Samples) == 0 {
        s.eventUOW.Rollback()
        return common.Failure[event.Event](&common.Error{Code: "400", Message: "Event samples are required"})
    }

    if len(req.EventAddress) == 0 {
        s.eventUOW.Rollback()
        return common.Failure[event.Event](&common.Error{Code: "400", Message: "Event address is required"})
    }

	// Create Event
	eventEntity := entities.Event{
		Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
		Name: req.Name,
		DateOfEvent: req.DateOfEvent,
		StartTime: req.StartTime,
		EndTime: req.EndTime,
		Limit: req.Limit,
		NumberSamples: req.NumberSamples,
		PhoneContact: req.PhoneContact,
		EmailContact: req.EmailContact,
		IsPublic: req.IsPublic,
		UserID: userID,
		RegisterDate: req.RegisterDate,
		RegisterStatus: eventregisterstatus.RegisterStatusEnumPending,
	}

    if err := eventRepo.Create(ctx, &eventEntity); err != nil {
        s.eventUOW.Rollback()
        return common.Failure[event.Event](&common.Error{Code: "500", Message: "Failed to create event"})
    }

	var sampleEntities []entities.UserSample
	var eventSampleEntities []entities.EventSample

	for _, sample := range req.Samples {
		sampleEntity := entities.UserSample{
			Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
			Name: sample.Name,
			RoastingDate: sample.RoastingDate,
			RoastLevel: sample.RoastLevel,
			AltitudeGrow: sample.AltitudeGrow,
			RoasteryName: sample.RoasteryName,
			RoasteryAddress: sample.RoasteryAddress,
			BreedName: sample.BreedName,
			PreProcessing: sample.PreProcessing,
			GrowNation: sample.GrowNation,
			GrowAddress: sample.GrowAddress,
			Price: sample.Price,
			UserID: userID,
		}
		sampleEntities = append(sampleEntities, sampleEntity)

		eventSampleEntity := entities.EventSample{
			Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
			Price: nil,
			Rating: nil,
			UserSampleID: sampleEntity.ID,
			EventID: eventEntity.ID,
		}
		eventSampleEntities = append(eventSampleEntities, eventSampleEntity)
	}

    if err := sampleRepo.CreateRange(ctx, sampleEntities); err != nil {
        s.eventUOW.Rollback()
        return common.Failure[event.Event](&common.Error{Code: "500", Message: "Failed to create event samples"})
    }

    if err := eventSampleRepo.CreateRange(ctx, eventSampleEntities); err != nil {
        s.eventUOW.Rollback()
        return common.Failure[event.Event](&common.Error{Code: "500", Message: "Failed to create event samples"})
    }

	var eventAddressesEntities []entities.EventAddress

	for _, eventAddress := range req.EventAddress {
		eventAddressEntity := entities.EventAddress{
			Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
			Province: eventAddress.Province,
			District: eventAddress.District,
			Longitude: eventAddress.Longitude,
			Latitude: eventAddress.Latitude,
			Ward: eventAddress.Ward,
			Street: eventAddress.Street,
			Phone: eventAddress.Phone,
			EventID: eventEntity.ID,
		}
		eventAddressesEntities = append(eventAddressesEntities, eventAddressEntity)
	}

    if err := eventAddressRepo.CreateRange(ctx, eventAddressesEntities); err != nil {
        s.eventUOW.Rollback()
        return common.Failure[event.Event](&common.Error{Code: "500", Message: "Failed to create event addresses"})
    }

	// Commit the transaction
    if err := s.eventUOW.Commit(); err != nil {
        s.eventUOW.Rollback()
        return common.Failure[event.Event](&common.Error{Code: "500", Message: "Failed to commit transaction"})
    }

    mapped := mapper.MapEventToContractGetAllEventResponse(&eventEntity)
    return common.Success(&mapped)
}

func (s *eventService) GetAll(ctx context.Context, pageSize int, pageNumber int, searchTerm string) (*common.Result[common.PageResult[event.Event]]) {
	db := s.eventUOW.GetDB()

    // Build (but do not execute) the query
    q := db.WithContext(ctx).Model(&entities.Event{})

    if searchTerm != "" {
        q = q.Where("name ILIKE ?", "%"+searchTerm+"%")
    }

	q = q.Where("is_public = ?", true)
	q = q.Where("is_deleted = ?", false)
	
    q = q.Order("created_at DESC")

	events, err := postgres.GetPaginated[entities.Event](q, ctx, pageSize, pageNumber, "EventAddress", "HostBy")
	if err != nil {
		return common.Failure[common.PageResult[event.Event]](&common.Error{Code: "500", Message: "Failed to get events"})
	}

	var eventsPageResult common.PageResult[event.Event]
	eventsPageResult.Data = make([]event.Event, len(events.Data))
	for i, event := range events.Data {
		eventsPageResult.Data[i] = mapper.MapEventToContractGetAllEventResponse(&event)
	}
	eventsPageResult.Total = int(events.Total)
	eventsPageResult.Page = events.Page
	eventsPageResult.PageSize = events.PageSize
	eventsPageResult.TotalPages = int(events.TotalPages)

	return common.Success(&eventsPageResult)
}

func (s *eventService) GetByID(ctx context.Context, id uuid.UUID) (*common.Result[event.GetEventByIDResponse]) {
	eventRepo := s.eventRepo
	// Preload all necessary relationships including nested UserSample
	eventEntity, err := eventRepo.GetByID(ctx, id, "EventAddress", "HostBy", "EventSamples.UserSample")
	if err != nil {
		return common.Failure[event.GetEventByIDResponse](&common.Error{Code: "500", Message: "Failed to get event"})
	}
	
	// Validate that we have a valid event entity
	if eventEntity == nil {
		return common.Failure[event.GetEventByIDResponse](&common.Error{Code: "404", Message: "Event not found"})
	}
	
	eventContract := mapper.MapEventToContractGetEventByIDResponse(eventEntity)
	return common.Success(&eventContract)
}

func (s *eventService) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}