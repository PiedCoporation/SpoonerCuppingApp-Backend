//go:build wireinject

package user

import (
	"backend/config"
	"backend/internal/infrastructure/repository/postgres"
	userInterface "backend/internal/usecase/user"
	userImpl "backend/internal/usecase/user/implement"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewUserAuthService(
	cfg *config.Config,
	db *gorm.DB,
) userInterface.UserAuthService {
	wire.Build(
		postgres.NewUserRepo,
		postgres.NewRefreshTokenRepo,
		userImpl.NewUserAuthService,
	)
	return nil
}
