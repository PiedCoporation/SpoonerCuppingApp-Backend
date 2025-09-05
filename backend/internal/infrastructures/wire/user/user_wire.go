//go:build wireinject

package user

import (
	postgres2 "backend/internal/persistents/postgres"
	userImpl "backend/internal/usecases"
	userInterface "backend/internal/usecases/abstractions"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewUserAuthService(
	db *gorm.DB,
) userInterface.UserAuthService {
	wire.Build(
		postgres2.NewUserAuthUow,
		postgres2.NewUserRepo,
		postgres2.NewRefreshTokenRepo,
		userImpl.NewUserAuthService,
	)
	return nil
}
