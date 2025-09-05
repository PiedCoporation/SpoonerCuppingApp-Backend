package initializations

import (
	"backend/global"
	"backend/internal/infrastructures/cache/rolecache"
	"backend/internal/usecases/abstractions"
	"context"

	"go.uber.org/zap"
)

func NewRolesCache(roleService abstractions.RoleService) {
	roles, err := roleService.GetAll(context.Background())
	if err != nil {
		global.Logger.Fatal("Roles initializations failed", zap.Error(err))
	}

	rolecache.NewCache(roles)
}
