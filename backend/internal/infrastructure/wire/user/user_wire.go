//go:build wireinject

package user

import (
	"backend/internal/infrastructure/repository/postgres"
	"backend/internal/infrastructure/uow"
	userInterface "backend/internal/usecase/user"
	userImpl "backend/internal/usecase/user/implement"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewUserAuthService(
	db *gorm.DB,
) userInterface.UserAuthService {
	wire.Build(
		uow.NewUserAuthUow,
		postgres.NewUserRepo,
		postgres.NewRefreshTokenRepo,
		userImpl.NewUserAuthService,
	)
	return nil
}
