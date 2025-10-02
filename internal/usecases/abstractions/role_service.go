package abstractions

import (
	"backend/internal/domains/entities"
	"context"
)

type IRoleService interface {
	GetAll(ctx context.Context) ([]entities.Role, error)
}
