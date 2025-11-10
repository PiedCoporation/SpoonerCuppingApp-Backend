package usecases

import (
	"backend/global"
	"backend/internal/constants/enums/eventparticipant"
	"backend/internal/contracts/common"
	eventContractRequest "backend/internal/contracts/event/request"
	eventContractResponse "backend/internal/contracts/event/response"
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

func (s *eventService) GetEventByUserID(ctx context.Context, pageSize int, pageNumber int) (*common.Result[common.PageResult[eventContractResponse.Event]]) {
	userID, _ := ctx.Value("userID").(uuid.UUID)

	// Build a query on events joined with event_users for the given user
	q := s.eventUOW.GetDB().WithContext(ctx).
		Model(&entities.Event{}).
		Joins("JOIN event_users eu ON eu.event_id = events.id").
		Where("eu.user_id = ?", userID).
		Where("eu.is_deleted = ?", false).
		Where("events.is_deleted = ?", false).
		Order("eu.is_host DESC").
		Order("events.created_at DESC")

	// Paginate events and preload relations required for mapping
	pg, err := postgres.GetPaginated[entities.Event](q, ctx, pageSize, pageNumber, "EventAddress", "HostBy")
	if err != nil {
		return common.Failure[common.PageResult[eventContractResponse.Event]](&common.Error{Code: 500, Message: "Failed to get events"})
	}

	var eventsPageResult common.PageResult[eventContractResponse.Event]
	eventsPageResult.Data = make([]eventContractResponse.Event, len(pg.Data))
	for i, ev := range pg.Data {
		eventsPageResult.Data[i] = mapper.MapEventToContractGetAllEventResponse(&ev)
	}
	eventsPageResult.Total = int(pg.Total)
	eventsPageResult.Page = pg.Page
	eventsPageResult.PageSize = pg.PageSize
	eventsPageResult.TotalPages = int(pg.TotalPages)

	return common.Success(&eventsPageResult)
}

func (s *eventService) StartEvent(ctx context.Context, id uuid.UUID) (*common.Result[string]) {
	userID, _ := ctx.Value("userID").(uuid.UUID)

	global.Logger.Info("userID", zap.Any("userID", userID))
	
	
	repoProvider, err := s.eventUOW.Begin(ctx)
	if err != nil {
		return common.Failure[string](&common.Error{Code: 500, Message: "Failed to begin transaction"})
	}
	
	eventRepo := repoProvider.EventRepository()

	eventEntity, err := eventRepo.GetByID(ctx, id)
	if err != nil {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 404, Message: "Event not found"})
	}

	if eventEntity.UserID != userID {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 403, Message: "You are not the host of this event"})
	}

	if eventEntity.IsStart {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 400, Message: "Event already started"})
	}

	eventEntity.IsStart = true
	if err := eventRepo.Update(ctx, eventEntity.ID, map[string]any{
		"is_start": eventEntity.IsStart,
	}); err != nil {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 500, Message: "Failed to update event"})
	}

	if err := s.eventUOW.Commit(); err != nil {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 500, Message: "Failed to commit transaction"})
	}

	msg := "Event started"

	return common.Success(&msg)
}

func (s *eventService) Register(ctx context.Context, id uuid.UUID) (*common.Result[string]) {
	userID, _ := ctx.Value("userID").(uuid.UUID)

	global.Logger.Info("userID", zap.Any("userID", userID))

	repoProvider, err := s.eventUOW.Begin(ctx)
	if err != nil {
		return common.Failure[string](&common.Error{Code: 500, Message: "Failed to begin transaction"})
	}

	eventRepo := repoProvider.EventRepository()
	eventUserRepo := repoProvider.EventUserRepository()

	eventEntity, err := eventRepo.GetByID(ctx, id)
	if err != nil {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 404, Message: "Event not found"})
	}

	if eventEntity.RegisterStartTime.After(time.Now()) {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 400, Message: "Event is not start for register"})
	}

	if eventEntity.RegisterEndTime.Before(time.Now()) {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 400, Message: "Event is not end for register"})
	}

	if eventEntity.TotalJoined >= eventEntity.Limit {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 400, Message: "Event is full"})
	}
	
	query := fmt.Sprintf("user_id = '%s' AND event_id = '%s'", userID.String(), id.String())

	// Check if user is already registered for this event
	existingEventUser, err := eventUserRepo.GetSingle(ctx, query)
	if err != nil {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 500, Message: "Failed to get event user"})
	}
	if existingEventUser != nil {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 400, Message: "User already registered"})
	}

	eventUserEntity := entities.EventUser{
		Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
		UserID: userID,
		EventID: id,
		IsAccepted: eventEntity.AutoAccept,
		IsInvited: false,
		IsHost: false,
	}

	if err := eventUserRepo.Create(ctx, &eventUserEntity); err != nil {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 500, Message: "Failed to create event user"})
	}

	updateData := map[string]any{}

	if(eventEntity.AutoAccept) {
		eventEntity.TotalJoined++
		updateData["total_joined"] = eventEntity.TotalJoined
	}
	eventEntity.TotalRegistered++
	updateData["total_registered"] = eventEntity.TotalRegistered
	if err := eventRepo.Update(ctx, eventEntity.ID, updateData); err != nil {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 500, Message: "Failed to update event"})
	}	

	if err := s.eventUOW.Commit(); err != nil {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 500, Message: "Failed to commit transaction"})
	}
	
	msg := "Registration successful"	
	return common.Success(&msg)
}

func (s *eventService) Create(ctx context.Context, req eventContractRequest.CreateEventReq) (*common.Result[eventContractResponse.Event]) {
	// g, gCtx := errgroup.WithContext(ctx)
	userID, _ := ctx.Value("userID").(uuid.UUID)

	repoProvider, err := s.eventUOW.Begin(ctx)
    if err != nil {
        return common.Failure[eventContractResponse.Event](&common.Error{Code: 500, Message: "Failed to begin transaction"})
    }

	sampleRepo := repoProvider.SampleRepository()
	eventRepo := repoProvider.EventRepository()
	eventAddressRepo := repoProvider.EventAddressRepository()
	eventSampleRepo := repoProvider.EventSampleRepository()
	eventUserRepo := repoProvider.EventUserRepository()

    if len(req.Samples) == 0 {
        s.eventUOW.Rollback()
        return common.Failure[eventContractResponse.Event](&common.Error{Code: 400, Message: "Event samples are required"})
    }

    if len(req.EventAddress) == 0 {
        s.eventUOW.Rollback()
        return common.Failure[eventContractResponse.Event](&common.Error{Code: 400, Message: "Event address is required"})
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
	}

    if err := eventRepo.Create(ctx, &eventEntity); err != nil {
        s.eventUOW.Rollback()
        return common.Failure[eventContractResponse.Event](&common.Error{Code: 500, Message: "Failed to create event"})
    }

	if err := eventUserRepo.Create(ctx, &entities.EventUser{
		Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
		UserID: userID,
		EventID: eventEntity.ID,
		IsHost: true,
		IsAccepted: false,
		IsInvited: false,
	}); err != nil {
		s.eventUOW.Rollback()
		return common.Failure[eventContractResponse.Event](&common.Error{Code: 500, Message: "Failed to create event user"})
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
        return common.Failure[eventContractResponse.Event](&common.Error{Code: 500, Message: "Failed to create event samples"})
    }

    if err := eventSampleRepo.CreateRange(ctx, eventSampleEntities); err != nil {
        s.eventUOW.Rollback()
        return common.Failure[eventContractResponse.Event](&common.Error{Code: 500, Message: "Failed to create event samples"})
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
        return common.Failure[eventContractResponse.Event](&common.Error{Code: 500, Message: "Failed to create event addresses"})
    }

	// Commit the transaction
    if err := s.eventUOW.Commit(); err != nil {
        s.eventUOW.Rollback()
        return common.Failure[eventContractResponse.Event](&common.Error{Code: 500, Message: "Failed to commit transaction"})
    }

    mapped := mapper.MapEventToContractGetAllEventResponse(&eventEntity)
    return common.Success(&mapped)
}

func (s *eventService) GetAll(ctx context.Context, pageSize int, pageNumber int, searchTerm string) (*common.Result[common.PageResult[eventContractResponse.Event]]) {
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
		return common.Failure[common.PageResult[eventContractResponse.Event]](&common.Error{Code: 500, Message: "Failed to get events"})
	}

	var eventsPageResult common.PageResult[eventContractResponse.Event]
	eventsPageResult.Data = make([]eventContractResponse.Event, len(events.Data))
	for i, event := range events.Data {
		eventsPageResult.Data[i] = mapper.MapEventToContractGetAllEventResponse(&event)
	}
	eventsPageResult.Total = int(events.Total)
	eventsPageResult.Page = events.Page
	eventsPageResult.PageSize = events.PageSize
	eventsPageResult.TotalPages = int(events.TotalPages)

	return common.Success(&eventsPageResult)
}

func (s *eventService) GetByID(ctx context.Context, id uuid.UUID) (*common.Result[eventContractResponse.GetEventByIDResponse]) {
	eventRepo := s.eventRepo
	// Preload all necessary relationships including nested UserSample
	eventEntity, err := eventRepo.GetByID(ctx, id, "EventAddress", "HostBy", "EventSamples.UserSample")
	if err != nil {
		return common.Failure[eventContractResponse.GetEventByIDResponse](&common.Error{Code: 500, Message: "Failed to get event"})
	}
	
	// Validate that we have a valid event entity
	if eventEntity == nil {
		return common.Failure[eventContractResponse.GetEventByIDResponse](&common.Error{Code: 404, Message: "Event not found"})
	}
	
	eventContract := mapper.MapEventToContractGetEventByIDResponse(eventEntity)
	return common.Success(&eventContract)
}

func (s *eventService) GetEventParticipant(ctx context.Context, eventID uuid.UUID, pageSize int, pageNumber int, searchTerm string, typeParticipant eventparticipant.TypeParticipantEnum) (*common.Result[common.PageResult[eventContractResponse.GetAllUserEventResponse]]) {
	db := s.eventUOW.GetDB()

	q := db.WithContext(ctx).
		Model(&entities.EventUser{}).
		Where("event_id = ?", eventID).
		Where("is_deleted = ?", false).
		Joins("JOIN users u ON u.id = event_users.user_id").
		Where("u.is_deleted = ?", false)

	switch typeParticipant {
	case eventparticipant.TypeParticipantEnumJoined:
		q = q.Where("event_users.is_accepted = ? AND event_users.is_host = ?", true, false)
	case eventparticipant.TypeParticipantEnumRegistered:
		// Requested (awaiting approval, not invited, not host)
		q = q.Where("event_users.is_accepted = ? AND event_users.is_invited = ? AND event_users.is_host = ?", false, false, false)
	case eventparticipant.TypeParticipantEnumInvited:
		q = q.Where("event_users.is_accepted = ? AND event_users.is_invited = ? AND event_users.is_host = ?", false, true, false)
	default:
		// No additional filter
	}

	if searchTerm != "" {
		like := "%" + searchTerm + "%"
		q = q.Where("(u.first_name ILIKE ? OR u.last_name ILIKE ? OR u.email ILIKE ? OR u.phone ILIKE ?)",
			like, like, like, like)
	}

	q = q.Order("event_users.created_at DESC")

	pg, err := postgres.GetPaginated[entities.EventUser](q, ctx, pageSize, pageNumber, "User")
	if err != nil {
		return common.Failure[common.PageResult[eventContractResponse.GetAllUserEventResponse]](&common.Error{Code: 500, Message: "Failed to get participants"})
	}

	var pageResult common.PageResult[eventContractResponse.GetAllUserEventResponse]
	pageResult.Data = make([]eventContractResponse.GetAllUserEventResponse, len(pg.Data))
	for i, eu := range pg.Data {
		pageResult.Data[i] = eventContractResponse.GetAllUserEventResponse{
			UserEventID: eu.ID,
			UserID:      eu.UserID,
			FirstName:   eu.User.FirstName,
			LastName:    eu.User.LastName,
			Email:       eu.User.Email,
			Phone:       eu.User.Phone,
			IsAccepted:  eu.IsAccepted,
			IsInvited:   eu.IsInvited,
			IsHost:      eu.IsHost,
		}
	}
	pageResult.Total = int(pg.Total)
	pageResult.Page = pg.Page
	pageResult.PageSize = pg.PageSize
	pageResult.TotalPages = int(pg.TotalPages)

	return common.Success(&pageResult)
}

func (s *eventService) ResponseEvent(ctx context.Context, eventUserID uuid.UUID, isAccept bool) (*common.Result[string]) {
	repoProvider, err := s.eventUOW.Begin(ctx)
	if err != nil {
		return common.Failure[string](&common.Error{Code: 500, Message: "Failed to begin transaction"})
	}

	eventRepo := repoProvider.EventRepository()
	eventUserRepo := repoProvider.EventUserRepository()

	// Load event user
	eventUser, err := eventUserRepo.GetByID(ctx, eventUserID)
	if err != nil || eventUser == nil {		
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 404, Message: "Event user not found"})
	}

	// Load event
	eventEntity, err := eventRepo.GetByID(ctx, eventUser.EventID)
	if err != nil || eventEntity == nil {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 404, Message: "Event not found"})
	}

	// Update acceptance
	if err := eventUserRepo.Update(ctx, eventUser.ID, map[string]any{
		"is_accepted": isAccept,
	}); err != nil {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 500, Message: "Failed to update event user"})
	}

	// If accepted, increase total joined
	if isAccept {
		eventEntity.TotalJoined++
		if err := eventRepo.Update(ctx, eventEntity.ID, map[string]any{
			"total_joined": eventEntity.TotalJoined,
		}); err != nil {
			s.eventUOW.Rollback()
			return common.Failure[string](&common.Error{Code: 500, Message: "Failed to update event totals"})
		}
	}

	if err := s.eventUOW.Commit(); err != nil {
		s.eventUOW.Rollback()
		return common.Failure[string](&common.Error{Code: 500, Message: "Failed to commit transaction"})
	}

	msg := "Response updated"
	return common.Success(&msg)
}

func (s *eventService) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}