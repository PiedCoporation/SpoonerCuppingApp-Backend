package sample

import (
	"backend/global"
	"backend/internal/constants/enums/jwtpurpose"
	wireSample "backend/internal/infrastructures/wire/sample"
	"backend/internal/presentations/http/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type SampleRouter struct {
}

func (r *SampleRouter) InitSampleRouter(
	router *gin.RouterGroup,
	db *gorm.DB,
) {
	cfg := global.Config

	eventController, _ := wireSample.InitSampleRouterHandler(db)

	eventGroup := router.Group("/samples")

	privateGroup := eventGroup.Group("")
	
	privateGroup.Use(middlewares.AuthHeader([]byte(cfg.JWT.AccessTokenKey), jwtpurpose.Access))
	{
		privateGroup.POST("", eventController.CreateSample)
		privateGroup.GET("", eventController.GetSamples)
	}	
}