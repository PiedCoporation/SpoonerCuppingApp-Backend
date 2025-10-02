package usecases

import (
	"backend/internal/domains/entities"
	repoAbstractions "backend/internal/persistents/abstractions"
	serviceAbstractions "backend/internal/usecases/abstractions"
	"context"
)

type roleService struct {
	roleRepo repoAbstractions.IRoleRepository
}

func NewRoleService(roleRepo repoAbstractions.IRoleRepository) serviceAbstractions.IRoleService {
	return &roleService{roleRepo: roleRepo}
}

// GetAll implements role.IRoleService.
func (rs *roleService) GetAll(ctx context.Context) ([]entities.Role, error) {
	roles, err := rs.roleRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}
