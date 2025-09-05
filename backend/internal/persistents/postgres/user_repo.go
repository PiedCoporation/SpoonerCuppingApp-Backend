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

func NewUserRepo(db *gorm.DB) abstractions.UserRepository {
	return &userPgRepo{
		genericRepository: NewGenericRepository[entities.User](db),
	}
}

// GetByEmail implements abstractions.UserRepository.
func (ur *userPgRepo) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := ur.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// IsPhoneTaken implements abstractions.UserRepository.
func (ur *userPgRepo) IsPhoneTaken(ctx context.Context, phone string, excludeUserID uuid.UUID) (bool, error) {
	var user entities.User
	err := ur.db.WithContext(ctx).
		Where("phone = ? AND id != ?", phone, excludeUserID).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	if user.Entity.IsDeleted {
		return true, errorcode.ErrPhoneBelongsToDeletedAccount
	}

	return true, nil
}

// IsEmailTaken implements abstractions.UserRepository.
func (ur *userPgRepo) IsEmailTaken(ctx context.Context, email string, excludeUserID uuid.UUID) (bool, error) {
	var user entities.User
	err := ur.db.WithContext(ctx).
		Where("email = ? AND id != ?", email, excludeUserID).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	if user.Entity.IsDeleted {
		return true, errorcode.ErrEmailBelongsToDeletedAccount
	}

	return true, nil
}
