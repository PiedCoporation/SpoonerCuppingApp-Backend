package postgres

import (
	"backend/internal/contracts/common"
	"context"
	"math"

	"gorm.io/gorm"
)

// GetPaginated retrieves paginated results for any entity type
func GetPaginated[T any](
	db *gorm.DB,
	ctx context.Context, pageSize int, pageNumber int, preloads ...string) (*common.PageResult[T], error) {
	// Validate pagination parameters
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var entities []T
	var total int64

	// Create base query
	db = db.WithContext(ctx)
	
	// Apply preloads
	for _, p := range preloads {
		db = db.Preload(p)
	}

	// Count total records
	if err := db.Model(new(T)).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (pageNumber - 1) * pageSize
	// Apply pagination and fetch data
	if err := db.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &common.PageResult[T]{
		Data:       entities,
		Total:      int(total),
		Page:       pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
