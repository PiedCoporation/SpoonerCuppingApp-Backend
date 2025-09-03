package repository

import (
	"backend/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	GenericRepository[entities.User]
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	IsPhoneTaken(ctx context.Context, phone string, excludeUserID uuid.UUID) (bool, error)
	IsEmailTaken(ctx context.Context, email string, excludeUserID uuid.UUID) (bool, error)
}
