package repository

import (
	"backend/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	UpdateEmailVerified(ctx context.Context, userID uuid.UUID, isVerified bool) error
	IsPhoneTaken(ctx context.Context, phone string, excludeUserID uuid.UUID) (bool, error)
	IsEmailTaken(ctx context.Context, email string, excludeUserID uuid.UUID) (bool, error)
}
