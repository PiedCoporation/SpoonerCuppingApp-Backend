package abstractions

import (
	"backend/internal/domains/entities"
	"context"

	"github.com/google/uuid"
)

type IRefreshTokenRepository interface {
	IGenericRepository[entities.RefreshToken]
	GetByTokenAndUserID(ctx context.Context, token string, userID uuid.UUID) (*entities.RefreshToken, error)
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
}
