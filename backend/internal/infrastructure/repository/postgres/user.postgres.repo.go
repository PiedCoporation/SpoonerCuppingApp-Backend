package postgres

import (
	"backend/internal/constants/errorcode"
	"backend/internal/domain/entities"
	"backend/internal/usecase/repository"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userPgRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repository.UserRepository {
	return &userPgRepo{db: db}
}

// GetByID implements repository.UserRepository.
func (ur *userPgRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := ur.db.WithContext(ctx).
		Preload("Role").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.ErrUserNotFound
		}
		return nil, err
	}

	if user.Entity.IsDeleted {
		return nil, errorcode.ErrDeletedAccount
	}

	return &user, nil
}

// GetByEmail implements repository.UserRepository.
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

	if user.Entity.IsDeleted {
		return nil, errorcode.ErrDeletedAccount
	}
	return &user, nil
}

// Create implements repository.UserRepository.
func (ur *userPgRepo) Create(ctx context.Context, user *entities.User) error {
	err := ur.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateEmailVerified implements repository.UserRepository.
func (ur *userPgRepo) UpdateEmailVerified(ctx context.Context, userID uuid.UUID, isVerified bool) error {
	if err := ur.db.WithContext(ctx).Model(&entities.User{}).
		Where("id = ?", userID).
		Update("is_verified", isVerified).Error; err != nil {
		return err
	}
	return nil
}

// IsPhoneTaken implements repository.UserRepository.
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

// IsEmailTaken implements repository.UserRepository.
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
