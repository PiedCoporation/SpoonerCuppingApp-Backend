package controller

import (
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/event"
	abstractions "backend/internal/usecases/abstractions"
	"backend/pkg/utils/validation"
	"context"
	"net/http"

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

	// Get userID from Gin context (set by auth middleware)
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