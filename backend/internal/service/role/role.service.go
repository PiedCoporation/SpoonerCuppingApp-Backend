package role

import (
	"backend/internal/domain/entities"
	"context"
)

type RoleService interface {
	GetAll(ctx context.Context) ([]entities.Role, error)
	GetByName(ctx context.Context, roleName string) (*entities.Role, error)
}
