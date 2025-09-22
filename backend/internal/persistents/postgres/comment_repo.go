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

// GetDirectChildren implements abstractions.ICommentRepository.
func (r *commentPgRepo) GetDirectChildren(ctx context.Context, parentID uuid.UUID) ([]entities.Comment, error) {
	return r.FindByQuery(ctx, "parent_id = ?", []any{parentID}, "User")
}

// GetRootComments implements abstractions.ICommentRepository.
func (r *commentPgRepo) GetRootComments(ctx context.Context, postID uuid.UUID) ([]entities.Comment, error) {
	return r.FindByQuery(ctx, "parent_id IS NULL AND post_id = ?", []any{postID}, "User")
}

// DeleteByPostID implements abstractions.ICommentRepository.
func (r *commentPgRepo) DeleteByPostID(ctx context.Context, postID uuid.UUID) error {
	var entity entities.Comment
	if err := r.db.WithContext(ctx).
		Model(&entity).
		Where("post_id = ?", postID).
		Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}
