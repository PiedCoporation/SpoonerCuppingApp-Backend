package usecases

import (
	"backend/internal/domains/entities"
	repoAbstractions "backend/internal/persistents/abstractions"
	serviceAbstractions "backend/internal/usecases/abstractions"
	"context"
)

type roleService struct {
	roleRepo repoAbstractions.RoleRepository
}

func NewRoleService(roleRepo repoAbstractions.RoleRepository) serviceAbstractions.RoleService {
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
