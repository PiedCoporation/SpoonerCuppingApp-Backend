package repository

import (
	"backend/internal/domain/entities"
	"context"
)

type RoleRepository interface {
	GetAll(ctx context.Context) ([]entities.Role, error)
	GetByName(ctx context.Context, name string) (*entities.Role, error)
}
