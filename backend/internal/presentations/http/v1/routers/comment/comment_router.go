package comment

import (
	"backend/global"
	"backend/internal/constants/enums/jwtpurpose"
	wireComment "backend/internal/infrastructures/wire/comment"
	"backend/internal/presentations/http/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentRouter struct{}

func (r *CommentRouter) InitCommentRouter(
	router *gin.RouterGroup,
	db *gorm.DB,
) {
	cfg := global.Config

	commentController, _ := wireComment.InitCommentRouterHandler(db)

	// ====== Main comment group ======
	commentGroup := router.Group("/comments")

	// ====== Public group ======
	publicGroup := commentGroup.Group("")
	{
		publicGroup.GET("/:id/children", commentController.GetDirectChildren)
	}

	// ====== Private group (using access token) ======
	privateGroup := commentGroup.Group("")
	privateGroup.Use(middlewares.AuthHeader([]byte(cfg.JWT.AccessTokenKey), jwtpurpose.Access))
	{
		privateGroup.POST("/", commentController.CreateComment)
		privateGroup.PUT("/:id", commentController.UpdateComment)
		privateGroup.DELETE("/:id", commentController.DeleteComment)
	}
}
