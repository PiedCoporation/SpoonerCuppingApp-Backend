package abstractions

import (
	"backend/internal/domains/entities"
	"context"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	GenericRepository[entities.RefreshToken]
	GetByTokenAndUserID(ctx context.Context, token string, userID uuid.UUID) (*entities.RefreshToken, error)
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
}
