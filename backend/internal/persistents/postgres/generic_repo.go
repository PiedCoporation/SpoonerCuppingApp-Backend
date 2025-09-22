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

// GetAll implements abstractions.GenericRepository.
func (r *genericRepository[T]) GetAll(ctx context.Context, preloads ...string) ([]T, error) {
	var entities []T
	db := r.db.WithContext(ctx)
	for _, p := range preloads {
		db = db.Preload(p)
	}
	if err := db.Where("is_deleted = ?", false).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// GetByID implements abstractions.GenericRepository.
func (r *genericRepository[T]) GetByID(ctx context.Context, id uuid.UUID, preloads ...string) (*T, error) {
	var entity T
	db := r.db.WithContext(ctx)
	for _, p := range preloads {
		db = db.Preload(p)
	}
	err := db.Where("id = ? AND is_deleted = ?", id, false).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.ErrNotFound
		}
		return nil, err
	}

	return &entity, nil
}

// GetSingle implements abstractions.GenericRepository.
func (r *genericRepository[T]) GetSingle(ctx context.Context, query string, preloads ...string) (*T, error) {
	var entity T
	db := r.db.WithContext(ctx)
	for _, p := range preloads {
		db = db.Preload(p)
	}
	err := db.Where(query).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindByQuery implements abstractions.GenericRepository.
func (r *genericRepository[T]) FindByQuery(
	ctx context.Context,
	query string, args []any, preloads ...string,
) ([]T, error) {
	var entities []T
	db := r.db.WithContext(ctx)
	for _, p := range preloads {
		db = db.Preload(p)
	}

	if err := db.
		Where("is_deleted = ?", false).
		Where(query, args...).
		Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

// Create implements abstractions.GenericRepository.
func (r *genericRepository[T]) Create(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).
		Create(entity).Error; err != nil {
		return err
	}
	return nil
}

// CreateRange implements abstractions.GenericRepository.
func (r *genericRepository[T]) CreateRange(ctx context.Context, entities []T) error {
	if err := r.db.WithContext(ctx).
		Create(&entities).Error; err != nil {
		return err
	}
	return nil
}

// Update implements abstractions.GenericRepository.
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

// SoftDelete implements abstractions.GenericRepository.
func (r *genericRepository[T]) SoftDelete(ctx context.Context, id uuid.UUID) error {
	var entity T
	if err := r.db.WithContext(ctx).
		Model(&entity).
		Where("id = ?", id).
		Update("is_deleted", true).Error; err != nil {
		return err
	}

	return nil
}

// Delete implements abstractions.GenericRepository.
func (r *genericRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	var entity T
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity).Error; err != nil {
		return err
	}
	return nil
}
