//go:build wireinject

package role

import (
	"backend/internal/repository/postgres"
	"backend/internal/service/role"
	roleImpl "backend/internal/service/role/implement"

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
