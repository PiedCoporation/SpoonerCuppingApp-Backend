package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type refreshTokenPgRepo struct {
	*genericRepository[entities.RefreshToken]
}

func NewRefreshTokenRepo(db *gorm.DB) abstractions.IRefreshTokenRepository {
	return &refreshTokenPgRepo{
		genericRepository: NewGenericRepository[entities.RefreshToken](db),
	}
}

func (r *refreshTokenPgRepo) GetByTokenAndUserID(ctx context.Context, token string, userID uuid.UUID) (*entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND token = ? AND revoked = false", userID, token).
		First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *refreshTokenPgRepo) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).Model(&entities.RefreshToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
	if err != nil {
		return err
	}
	return nil
}
