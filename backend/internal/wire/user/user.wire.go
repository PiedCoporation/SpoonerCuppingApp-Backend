//go:build wireinject

package user

import (
	"backend/config"
	"backend/internal/repository/postgres"
	userInterface "backend/internal/service/user"
	userImpl "backend/internal/service/user/implement"

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
