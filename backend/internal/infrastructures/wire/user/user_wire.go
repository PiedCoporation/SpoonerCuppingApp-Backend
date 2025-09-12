//go:build wireinject

package user

import (
	postgres2 "backend/internal/persistents/postgres"
	"backend/internal/presentations/http/v1/controller"
	userImpl "backend/internal/usecases"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitUserRouterHandler(
	db *gorm.DB,
) (*controller.UserAuthController, error) {
	wire.Build(
		postgres2.NewUserAuthUow,
		postgres2.NewUserRepo,
		postgres2.NewRefreshTokenRepo,
		userImpl.NewUserAuthService,
		controller.NewUserAuthController,
	)
	return &controller.UserAuthController{}, nil
}
