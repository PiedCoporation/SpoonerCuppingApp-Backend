package http

import (
	"backend/config"
	"backend/internal/controller/http/v1/router"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	cfg *config.Config,
	userServiceSet *router.UserServiceSet,
) *gin.Engine {
	r := gin.Default()

	MainGroup := r.Group("/v1")
	{
		router.NewUserRouter(MainGroup, cfg, userServiceSet)
	}

	return r
}
