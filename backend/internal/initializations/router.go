package initializations

import (
	"backend/internal/presentations/http/v1/routers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	docs "backend/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	wsConfig "backend/internal/presentations/ws/configs"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	hub := wsConfig.NewHub()
	go hub.InitialHub()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "coffee-cupping-backend",
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		// TODO: validate token/JWT/etc if needed
		ServeWS(hub, c.Writer, c.Request)
	})

	// Swagger docs configuration
	docs.SwaggerInfo.BasePath = "/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	eventRouter := routers.RouterGroupApp.Event
	userRouter := routers.RouterGroupApp.User
	postRouter := routers.RouterGroupApp.Post

	MainGroup := r.Group("/v1")
	{
		eventRouter.InitEventRouter(MainGroup, db)
		userRouter.InitUserRouter(MainGroup, db)
		postRouter.InitPostRouter(MainGroup, db)
	}

	return r
}
