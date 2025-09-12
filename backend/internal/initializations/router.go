package initializations

import (
	"backend/internal/presentations/http/v1/routers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	eventRouter := routers.RouterGroupApp.Event
	userRouter := routers.RouterGroupApp.User

	MainGroup := r.Group("/v1")
	{
		eventRouter.InitEventRouter(MainGroup, db)
		userRouter.InitUserRouter(MainGroup, db)
	}

	return r
}