package event

import (
	"backend/internal/constants/enums/jwtpurpose"
	wireEvent "backend/internal/infrastructures/wire/event"
	"backend/internal/presentations/http/middlewares"

	"backend/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EventRouter struct {
}

func (u *EventRouter) InitEventRouter(
	router *gin.RouterGroup,
	db *gorm.DB,
) {
	cfg := global.Config

	eventController, _ := wireEvent.InitEventRouterHandler(db)

	eventGroup := router.Group("/events")

	privateGroup := eventGroup.Group("")
	
	privateGroup.Use(middlewares.AuthHeader([]byte(cfg.JWT.AccessTokenKey), jwtpurpose.Access))
	{
		privateGroup.POST("/", eventController.CreateEvent)
		privateGroup.GET("/", eventController.GetEvents)
		privateGroup.GET("/:id", eventController.GetEventByID)
		privateGroup.POST("/:id/register", eventController.RegisterEvent)
	}
}