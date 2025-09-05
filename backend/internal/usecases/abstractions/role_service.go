package abstractions

import (
	"backend/internal/domains/entities"
	"context"
)

type RoleService interface {
	GetAll(ctx context.Context) ([]entities.Role, error)
}
