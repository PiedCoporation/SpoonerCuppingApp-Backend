package initialization

import (
	"backend/internal/infrastructure/cache/rolecache"
	"backend/internal/usecase/role"
	"backend/pkg/logger"
	"context"

	"go.uber.org/zap"
)

func NewRolesCache(roleService role.RoleService, logger logger.Interface) {
	roles, err := roleService.GetAll(context.Background())
	if err != nil {
		logger.Fatal("Roles initialization failed", zap.Error(err))
	}

	rolecache.NewCache(roles)
}
