package initialization

import (
	"backend/global"
	"backend/internal/infrastructure/cache/rolecache"
	"backend/internal/usecase/role"
	"context"

	"go.uber.org/zap"
)

func NewRolesCache(roleService role.RoleService) {
	roles, err := roleService.GetAll(context.Background())
	if err != nil {
		global.Logger.Fatal("Roles initialization failed", zap.Error(err))
	}

	rolecache.NewCache(roles)
}
