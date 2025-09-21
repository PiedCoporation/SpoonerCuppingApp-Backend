package post

import (
	"backend/global"
	"backend/internal/constants/enums/jwtpurpose"
	wirePost "backend/internal/infrastructures/wire/post"
	"backend/internal/presentations/http/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostRouter struct{}

func (u *PostRouter) InitPostRouter(
	router *gin.RouterGroup,
	db *gorm.DB,
) {
	cfg := global.Config

	postController, _ := wirePost.InitPostRouterHandler(db)
	// ====== Main post group ======
	postGroup := router.Group("/posts")

	// ====== Public group ======
	publicGroup := postGroup.Group("")
	{
		publicGroup.GET("/", postController.GetPosts)
		publicGroup.GET("/:id", postController.GetPostByID)
	}

	// ====== Private group (using access token) ======
	privateGroup := postGroup.Group("")
	privateGroup.Use(middlewares.AuthHeader([]byte(cfg.JWT.AccessTokenKey), jwtpurpose.Access))
	{
		privateGroup.POST("/", postController.CreatePost)
		privateGroup.PATCH("/:id", postController.UpdatePost)
		privateGroup.DELETE("/:id", postController.DeletePost)
	}
}
