package controller

import (
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/event"
	abstractions "backend/internal/usecases/abstractions"
	"backend/pkg/utils/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	EventService abstractions.IEventService
}

func NewEventController(eventService abstractions.IEventService) *EventController {
	return &EventController{
		EventService: eventService,
	}
}

func (ec *EventController) CreateEvent(c *gin.Context) {
	var req event.CreateEventReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	ctx := c.Request.Context()
	err := ec.EventService.Create(ctx, req)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}
}