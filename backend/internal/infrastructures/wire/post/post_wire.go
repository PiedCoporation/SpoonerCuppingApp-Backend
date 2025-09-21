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
		postgres.NewCommentRepo,
		postgres.NewPostUow,
		usecases.NewPostService,
		controller.NewPostController,
	)
	return &controller.PostController{}, nil
}
