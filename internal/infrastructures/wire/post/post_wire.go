//go:build wireinject

//go:generate wire
package post

import (
	"backend/internal/persistents/postgres"
	"backend/internal/presentations/http/v1/controller"
	"backend/internal/usecases"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitPostRouterHandler(
	db *gorm.DB,
) (*controller.PostController, error) {
	wire.Build(
		postgres.NewPostRepo,
		postgres.NewPostImageRepo,
		postgres.NewPostLikeRepo,
		postgres.NewPostCommentRepo,
		postgres.NewPostUow,
		usecases.NewPostService,
		usecases.NewPostLikeService,
		controller.NewPostController,
	)
	return &controller.PostController{}, nil
}

func InitPostCommentRouterHandler(
	db *gorm.DB,
) (*controller.PostCommentController, error) {
	wire.Build(
		postgres.NewPostCommentRepo,
		postgres.NewPostRepo,
		usecases.NewPostCommentService,
		controller.NewPostCommentController,
	)
	return &controller.PostCommentController{}, nil
}
