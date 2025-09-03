package role

import (
	"backend/internal/domain/entities"
	"context"
)

type RoleService interface {
	GetAll(ctx context.Context) ([]entities.Role, error)
}
