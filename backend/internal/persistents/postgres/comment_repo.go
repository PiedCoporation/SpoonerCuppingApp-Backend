package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type commentPgRepo struct {
	*genericRepository[entities.Comment]
}

func NewCommentRepo(db *gorm.DB) abstractions.ICommentRepository {
	return &commentPgRepo{
		genericRepository: NewGenericRepository[entities.Comment](db),
	}
}

// DeleteByPostID implements abstractions.ICommentRepository.
func (p *commentPgRepo) DeleteByPostID(ctx context.Context, postID uuid.UUID) error {
	var entity entities.Comment
	if err := p.db.WithContext(ctx).
		Model(&entity).
		Where("post_id = ?", postID).
		Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}
