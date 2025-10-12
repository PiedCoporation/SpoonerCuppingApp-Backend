package controller

import (
	"backend/internal/contracts/sample"
	abstractions "backend/internal/usecases/abstractions"
	"backend/pkg/utils/validation"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SampleController struct {
	sampleService     abstractions.ISampleService
}

func NewSampleController(
	sampleService abstractions.ISampleService,
) *SampleController {
	return &SampleController{
		sampleService: sampleService,
	}
}

func (sc *SampleController) GetSamples(c *gin.Context) {
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil {
		pageSize = 10
	}
	pageNumber, err := strconv.Atoi(c.Query("page_number"))
	if err != nil {
		pageNumber = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if pageNumber < 1 {
		pageNumber = 1
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "userID", userID.(uuid.UUID))

	events := sc.sampleService.GetAll(ctx, pageSize, pageNumber)
	c.JSON(http.StatusOK, events)
}

func (sc *SampleController) CreateSample(c *gin.Context) {
	var req sample.SampleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}
	
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "userID", userID.(uuid.UUID))

	sample := sc.sampleService.Create(ctx, req)
	c.JSON(http.StatusCreated, sample)
}