package repository

import (
	"backend/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	GenericRepository[entities.RefreshToken]
	GetByTokenAndUserID(ctx context.Context, token string, userID uuid.UUID) (*entities.RefreshToken, error)
	// Revoke(ctx context.Context, id uuid.UUID) error
}
