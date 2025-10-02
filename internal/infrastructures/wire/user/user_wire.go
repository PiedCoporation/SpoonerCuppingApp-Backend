//go:build wireinject

//go:generate wire

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
) (*controller.UserController, error) {
	wire.Build(
		postgres2.NewUserRepo,
		userImpl.NewUserSettingService,
		controller.NewUserController,
	)
	return &controller.UserController{}, nil
}

func InitUserAuthRouterHandler(
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
