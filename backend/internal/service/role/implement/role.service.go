package implement

import (
	"backend/internal/domain/entities"
	"backend/internal/repository"
	"backend/internal/service/role"
	"context"
)

type roleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) role.RoleService {
	return &roleService{roleRepo: roleRepo}
}

// GetAll implements role.RoleService.
func (rs *roleService) GetAll(ctx context.Context) ([]entities.Role, error) {
	roles, err := rs.roleRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetByName implements role.RoleService.
func (rs *roleService) GetByName(ctx context.Context, roleName string) (*entities.Role, error) {
	role, err := rs.roleRepo.GetByName(ctx, roleName)
	if err != nil {
		return nil, err
	}
	return role, nil
}
