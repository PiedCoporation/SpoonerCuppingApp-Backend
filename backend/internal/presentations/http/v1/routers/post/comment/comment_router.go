package comment

import (
	"backend/global"
	"backend/internal/constants/enums/jwtpurpose"
	wirePost "backend/internal/infrastructures/wire/post"
	"backend/internal/presentations/http/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostCommentRouter struct{}

func (r *PostCommentRouter) InitPostCommentRouter(
	postRouter *gin.RouterGroup,
	db *gorm.DB,
) {
	cfg := global.Config

	commentController, _ := wirePost.InitPostCommentRouterHandler(db)

	// ====== Main postcomment group ======
	commentGroup := postRouter.Group("/:id/comments")

	// ====== Public group ======
	publicGroup := commentGroup.Group("")
	{
		publicGroup.GET("/", commentController.GetRootComments)
		publicGroup.GET("/:commentId/children", commentController.GetDirectChildren)
	}

	// ====== Private group (using access token) ======
	privateGroup := commentGroup.Group("")
	privateGroup.Use(middlewares.AuthHeader([]byte(cfg.JWT.AccessTokenKey), jwtpurpose.Access))
	{
		privateGroup.POST("/", commentController.CreateComment)
		privateGroup.PUT("/:commentId", commentController.UpdateComment)
		privateGroup.DELETE("/:commentId", commentController.DeleteComment)
	}
}
