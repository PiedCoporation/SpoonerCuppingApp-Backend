package postgres

import (
	"backend/internal/constants/errorcode"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type genericRepository[T any] struct {
	db *gorm.DB
}

func NewGenericRepository[T any](
	db *gorm.DB,
) *genericRepository[T] {
	return &genericRepository[T]{
		db: db,
	}
}

// GetAll implements repository.GenericRepository.
func (r *genericRepository[T]) GetAll(ctx context.Context, preloads ...string) ([]T, error) {
	var entities []T
	db := r.db.WithContext(ctx)
	for _, p := range preloads {
		db = db.Preload(p)
	}
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// GetByID implements repository.GenericRepository.
func (r *genericRepository[T]) GetByID(ctx context.Context, id uuid.UUID, preloads ...string) (*T, error) {
	var entity T
	db := r.db.WithContext(ctx)
	for _, p := range preloads {
		db = db.Preload(p)
	}
	err := db.Where("id = ?", id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.ErrNotFound
		}
		return nil, err
	}

	return &entity, nil
}

// Create implements repository.GenericRepository.
func (r *genericRepository[T]) Create(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).
		Create(entity).Error; err != nil {
		return err
	}
	return nil
}

// Update implements repository.GenericRepository.
func (r *genericRepository[T]) Update(ctx context.Context, id uuid.UUID, fields map[string]any) error {
	var entity T
	if err := r.db.WithContext(ctx).
		Model(&entity).
		Where("id = ?", id).
		Updates(fields).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements repository.GenericRepository.
func (r *genericRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	var entity T
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity).Error; err != nil {
		return err
	}
	return nil
}
