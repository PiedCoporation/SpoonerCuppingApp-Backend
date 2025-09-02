package repository

import (
	"backend/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	GetByTokenAndUserID(ctx context.Context, token string, userID uuid.UUID) (*entities.RefreshToken, error)
	Create(ctx context.Context, refreshToken *entities.RefreshToken) error
	Revoke(ctx context.Context, token string) error
}
