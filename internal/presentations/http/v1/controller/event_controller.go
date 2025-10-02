package controller

import (
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/event"
	abstractions "backend/internal/usecases/abstractions"
	"backend/pkg/utils/validation"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventController struct {
	EventService abstractions.IEventService
}

func NewEventController(eventService abstractions.IEventService) *EventController {
	return &EventController{
		EventService: eventService,
	}
}

// CreateEvent godoc
// @Summary Create an event
// @Description Create a new cupping event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body event.CreateEventReq true "Event payload"
// @Success 201 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse
// @Failure 401 {object} controller.ErrorResponse
// @Failure 500 {object} controller.ErrorResponse
// @Router /events [post]
func (ec *EventController) CreateEvent(c *gin.Context) {
	var req event.CreateEventReq
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
	err := ec.EventService.Create(ctx, req)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event created"})
}

// GetEvents godoc
// @Summary Get paginated events
// @Description Retrieve a paginated list of events with optional search functionality
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page_size query int false "Number of events per page (default: 10, max: 100)" default(10) minimum(1) maximum(100)
// @Param page_number query int false "Page number (default: 1)" default(1) minimum(1)
// @Param search_term query string false "Search term to filter events by name"
// @Success 200 {object} event.EventPageResult "Paginated list of events"
// @Failure 400 {object} controller.ErrorResponse "Bad request - invalid parameters"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized - invalid or missing token"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /events [get]
func (ec *EventController) GetEvents(c *gin.Context) {
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

	searchTerm := c.Query("search_term")
	ctx := c.Request.Context()
	events, err := ec.EventService.GetAll(ctx, pageSize, pageNumber, searchTerm)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEventByID godoc
// @Summary Get event by ID
// @Description Retrieve a specific event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID" format(uuid)
// @Success 200 {object} event.GetEventByIDResponse "Event details with samples"
// @Failure 400 {object} controller.ErrorResponse "Bad request - invalid event ID"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized - invalid or missing token"
// @Failure 404 {object} controller.ErrorResponse "Event not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /events/{id} [get]
func (ec *EventController) GetEventByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}
	ctx := c.Request.Context()

	event, err := ec.EventService.GetByID(ctx, id)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}
	c.JSON(http.StatusOK, event)
}

// RegisterEvent godoc
// @Summary Register for an event
// @Description Register the current user for a specific event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID" format(uuid)
// @Success 200 {object} controller.MessageResponse "Registration successful"
// @Failure 400 {object} controller.ErrorResponse "Bad request - invalid event ID or already registered"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized - invalid or missing token"
// @Failure 404 {object} controller.ErrorResponse "Event not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /events/{id}/register [post]
func (ec *EventController) RegisterEvent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}
	
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "userID", userID.(uuid.UUID))

	err = ec.EventService.Register(ctx, id)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}