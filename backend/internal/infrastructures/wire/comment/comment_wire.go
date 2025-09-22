//go:build wireinject

//go:generate wire
package comment

import (
	"backend/internal/persistents/postgres"
	"backend/internal/presentations/http/v1/controller"
	"backend/internal/usecases"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitCommentRouterHandler(
	db *gorm.DB,
) (*controller.CommentController, error) {
	wire.Build(
		postgres.NewCommentRepo,
		postgres.NewPostRepo,
		usecases.NewCommentService,
		controller.NewCommentController,
	)
	return &controller.CommentController{}, nil
}
