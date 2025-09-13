package usecases

import (
	eventEnum "backend/internal/constants/enums/event"
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/common"
	"backend/internal/contracts/event"
	"backend/internal/domains/commons"
	"backend/internal/domains/entities"
	persistentRepo "backend/internal/persistents/abstractions"
	"backend/internal/persistents/postgres"
	abstractions "backend/internal/usecases/abstractions"
	"context"

	"github.com/google/uuid"
)

type eventService struct {
	eventUOW persistentRepo.EventUOW
}

func NewEventService(eventUOW persistentRepo.EventUOW) abstractions.IEventService {
	return &eventService{
		eventUOW: eventUOW,
	}
}

func (s *eventService) Create(ctx context.Context, req event.CreateEventReq) error {
	// g, gCtx := errgroup.WithContext(ctx)
	userID, _ := ctx.Value("userID").(uuid.UUID)

	repoProvider, err := s.eventUOW.Begin(ctx)
	if err != nil {
		return err
	}

	sampleRepo := repoProvider.SampleRepository()
	eventRepo := repoProvider.EventRepository()
	eventAddressRepo := repoProvider.EventAddressRepository()
	eventSampleRepo := repoProvider.EventSampleRepository()

	if len(req.Samples) == 0 {
		s.eventUOW.Rollback()
		return errorcode.ErrEventSamplesRequired
	}

	if len(req.EventAddress) == 0 {
		s.eventUOW.Rollback()
		return errorcode.ErrEventAddressRequired
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
		RegisterStatus: eventEnum.RegisterStatusEnumPending,
	}

	if err := eventRepo.Create(ctx, &eventEntity); err != nil {
		s.eventUOW.Rollback()
		return err
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
		return err
	}

	if err := eventSampleRepo.CreateRange(ctx, eventSampleEntities); err != nil {
		s.eventUOW.Rollback()
		return err
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
		return err
	}

	// Commit the transaction
	if err := s.eventUOW.Commit(); err != nil {
		s.eventUOW.Rollback()
		return err
	}

	return nil
}

func (s *eventService) GetAll(ctx context.Context, pageSize int, pageNumber int, searchTerm string) (*common.PageResult[event.Event], error) {
	db := s.eventUOW.GetDB()

    // Build (but do not execute) the query
    q := db.WithContext(ctx).Model(&entities.Event{})

    // Optional search (use ILIKE for Postgres; switch to LIKE for MySQL)
    if searchTerm != "" {
        q = q.Where("name ILIKE ?", "%"+searchTerm+"%")
    }

    // Optional default ordering for stable pagination
    q = q.Order("created_at DESC")

	events, err := postgres.GetPaginated[entities.Event](q, ctx, pageSize, pageNumber, "EventAddress")
	if err != nil {
		return nil, err
	}

	var eventsPageResult common.PageResult[event.Event]
	eventsPageResult.Data = make([]event.Event, len(events.Data))
	for i, event := range events.Data {
		eventsPageResult.Data[i] = event.ToContract()
	}
	eventsPageResult.Total = int(events.Total)
	eventsPageResult.Page = events.Page
	eventsPageResult.PageSize = events.PageSize
	eventsPageResult.TotalPages = int(events.TotalPages)

	return &eventsPageResult, nil
}


func (s *eventService) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}