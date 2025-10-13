package usecases

import (
	"backend/internal/contracts/common"
	"backend/internal/contracts/sample"
	"backend/internal/domains/commons"
	"backend/internal/domains/entities"
	"backend/internal/mapper"
	repoAbstractions "backend/internal/persistents/abstractions"
	"backend/internal/persistents/postgres"
	abstractions "backend/internal/usecases/abstractions"
	"context"

	"github.com/google/uuid"
)

type sampleService struct {
	sampleRepo repoAbstractions.ISampleRepository
	sampleUOW repoAbstractions.SampleUOW
}

func NewSampleService(sampleRepo repoAbstractions.ISampleRepository, sampleUOW repoAbstractions.SampleUOW) abstractions.ISampleService {
	return &sampleService{sampleRepo: sampleRepo, sampleUOW: sampleUOW}
}

func (s *sampleService) Create(ctx context.Context, req sample.SampleReq) (*common.Result[sample.SampleRes]) {
	// Validate userID from context
	userID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return common.Failure[sample.SampleRes](&common.Error{Code: 401, Message: "Invalid user context"})
	}

	// Validate required fields
	if req.Name == "" {
		return common.Failure[sample.SampleRes](&common.Error{Code: 400, Message: "Sample name is required"})
	}
	if req.RoasteryName == "" {
		return common.Failure[sample.SampleRes](&common.Error{Code: 400, Message: "Roastery name is required"})
	}
	if req.BreedName == "" {
		return common.Failure[sample.SampleRes](&common.Error{Code: 400, Message: "Breed name is required"})
	}
	if req.Price <= 0 {
		return common.Failure[sample.SampleRes](&common.Error{Code: 400, Message: "Price must be greater than 0"})
	}


	repoProvider, err := s.sampleUOW.Begin(ctx)
	if err != nil {
		return common.Failure[sample.SampleRes](&common.Error{Code: 500, Message: "Failed to begin transaction"})
	}

	sampleRepo := repoProvider.SampleRepository()
	sampleTastingRepo := repoProvider.SampleTastingRepository()

	sampleEntity := entities.UserSample{
		Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
		Name: req.Name,
		RoastingDate: req.RoastingDate,
		RoastLevel: req.RoastLevel,
		AltitudeGrow: req.AltitudeGrow,
		RoasteryName: req.RoasteryName,
		RoasteryAddress: req.RoasteryAddress,
		BreedName: req.BreedName,
		PreProcessing: req.PreProcessing,
		GrowNation: req.GrowNation,
		GrowAddress: req.GrowAddress,
		Price: req.Price,
		UserID: userID,
	}

	// Create the sample first
	if err := sampleRepo.Create(ctx, &sampleEntity); err != nil {
		s.sampleUOW.Rollback()
		return common.Failure[sample.SampleRes](&common.Error{Code: 500, Message: "Failed to create sample"})
	}

	// Create the tasting records separately
	for _, tasting := range req.UserSampleTastings {
		tastingEntity := entities.UserSampleTasting{
			Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
			ParentName: tasting.ParentName,
			ChildName: tasting.ChildName,
			GrandChildName: tasting.GrandChildName,
			UserSampleID: sampleEntity.ID,
		}
		
		if err := sampleTastingRepo.Create(ctx, &tastingEntity); err != nil {
			s.sampleUOW.Rollback()
			return common.Failure[sample.SampleRes](&common.Error{Code: 500, Message: "Failed to create sample tasting"})
		}
	}

	if err := s.sampleUOW.Commit(); err != nil {
		return common.Failure[sample.SampleRes](&common.Error{Code: 500, Message: "Failed to commit transaction"})
	}

	// Use the created entity directly instead of fetching from database
	// We need to attach the tasting data to the sample entity for the response
	sampleEntity.UserSampleTastings = make([]entities.UserSampleTasting, len(req.UserSampleTastings))
	for i, tasting := range req.UserSampleTastings {
		sampleEntity.UserSampleTastings[i] = entities.UserSampleTasting{
			Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
			ParentName: tasting.ParentName,
			ChildName: tasting.ChildName,
			GrandChildName: tasting.GrandChildName,
			UserSampleID: sampleEntity.ID,
		}
	}

	mapped := mapper.MapUserSampleToContractGetAllSampleResponse(&sampleEntity)

	return common.Success(mapped)
}

func (s *sampleService) GetAll(ctx context.Context, pageSize int, pageNumber int) (*common.Result[common.PageResult[sample.SampleRes]]) {
	// Validate userID from context
	userID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return common.Failure[common.PageResult[sample.SampleRes]](&common.Error{Code: 401, Message: "Invalid user context"})
	}

	q := s.sampleUOW.GetDB().WithContext(ctx).
		Model(&entities.UserSample{}).
		Preload("UserSampleTastings").
		Where("user_id = ?", userID).
		Where("is_deleted = ?", false).
		Order("created_at DESC")

	pg, err := postgres.GetPaginated[entities.UserSample](q, ctx, pageSize, pageNumber)
	if err != nil {
		return common.Failure[common.PageResult[sample.SampleRes]](&common.Error{Code: 500, Message: "Failed to get samples"})
	}

	var samplesPageResult common.PageResult[sample.SampleRes]
	samplesPageResult.Data = make([]sample.SampleRes, len(pg.Data))
	for i, sample := range pg.Data {
		samplesPageResult.Data[i] = *mapper.MapUserSampleToContractGetAllSampleResponse(&sample)
	}
	samplesPageResult.Total = int(pg.Total)
	samplesPageResult.Page = pg.Page
	samplesPageResult.PageSize = pg.PageSize
	samplesPageResult.TotalPages = int(pg.TotalPages)

	return common.Success(&samplesPageResult)
}