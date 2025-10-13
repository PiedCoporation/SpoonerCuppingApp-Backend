package postgres

import (
	"backend/internal/constants/errorcode"
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userPgRepo struct {
	*genericRepository[entities.User]
}

func NewUserRepo(db *gorm.DB) abstractions.IUserRepository {
	return &userPgRepo{
		genericRepository: NewGenericRepository[entities.User](db),
	}
}

// GetByEmail implements abstractions.IUserRepository.
func (ur *userPgRepo) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := ur.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// IsPhoneTaken implements abstractions.IUserRepository.
func (ur *userPgRepo) IsPhoneTaken(ctx context.Context, phone string, excludeUserID uuid.UUID) (bool, error) {
	var count int64
	err := ur.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("phone = ? AND id != ? AND is_deleted = ?", phone, excludeUserID, false).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsEmailTaken implements abstractions.IUserRepository.
func (ur *userPgRepo) IsEmailTaken(ctx context.Context, email string, excludeUserID uuid.UUID) (bool, error) {
	var count int64
	err := ur.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("email = ? AND id != ? AND is_deleted = ?", email, excludeUserID, false).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
