//go:build wireinject

package role

import (
	"backend/internal/infrastructure/repository/postgres"
	"backend/internal/usecase/role"
	roleImpl "backend/internal/usecase/role/implement"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewRoleService(db *gorm.DB) role.RoleService {
	wire.Build(
		postgres.NewRoleRepo,
		roleImpl.NewRoleService,
	)
	return nil
}
