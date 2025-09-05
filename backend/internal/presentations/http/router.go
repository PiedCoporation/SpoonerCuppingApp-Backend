package http

import (
	"backend/internal/presentations/http/v1/router"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	userServiceSet *router.UserServiceSet,
) *gin.Engine {
	r := gin.Default()

	MainGroup := r.Group("/v1")
	{
		router.NewUserRouter(MainGroup, userServiceSet)
	}

	return r
}
