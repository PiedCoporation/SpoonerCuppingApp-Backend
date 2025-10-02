//go:build wireinject

//go:generate wire

package role

import (
	"backend/internal/persistents/postgres"
	roleImpl "backend/internal/usecases"
	"backend/internal/usecases/abstractions"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewRoleService(db *gorm.DB) abstractions.IRoleService {
	wire.Build(
		postgres.NewRoleRepo,
		roleImpl.NewRoleService,
	)
	return nil
}
