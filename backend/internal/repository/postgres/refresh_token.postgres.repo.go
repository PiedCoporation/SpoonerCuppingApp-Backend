package postgres

import (
	"backend/internal/domain/entities"
	"backend/internal/repository"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type refreshTokenPgRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepo(db *gorm.DB) repository.RefreshTokenRepository {
	return &refreshTokenPgRepo{db: db}
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

func (r *refreshTokenPgRepo) Create(ctx context.Context, refreshToken *entities.RefreshToken) error {
	err := r.db.WithContext(ctx).Create(&refreshToken).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *refreshTokenPgRepo) Revoke(ctx context.Context, token string) error {
	err := r.db.WithContext(ctx).Model(&entities.RefreshToken{}).
		Where("token = ? AND revoked = false", token).
		Update("revoked", true).Error
	if err != nil {
		return err
	}
	return nil
}
