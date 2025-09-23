package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postCommentPgRepo struct {
	*genericRepository[entities.PostComment]
}

func NewPostCommentRepo(db *gorm.DB) abstractions.IPostCommentRepository {
	return &postCommentPgRepo{
		genericRepository: NewGenericRepository[entities.PostComment](db),
	}
}

// GetDirectChildren implements abstractions.ICommentRepository.
func (r *postCommentPgRepo) GetDirectChildren(
	ctx context.Context,
	parentID uuid.UUID, orderByCreatedAtDesc bool,
) ([]entities.PostComment, error) {
	return r.FindByQuery(ctx, "parent_id = ?", []any{parentID}, orderByCreatedAtDesc, "User")
}

// GetRootComments implements abstractions.ICommentRepository.
func (r *postCommentPgRepo) GetRootComments(
	ctx context.Context,
	postID uuid.UUID, orderByCreatedAtDesc bool,
) ([]entities.PostComment, error) {
	return r.FindByQuery(ctx, "parent_id IS NULL AND post_id = ?", []any{postID}, orderByCreatedAtDesc, "User")
}

// DeleteByPostID implements abstractions.ICommentRepository.
func (r *postCommentPgRepo) DeleteByPostID(ctx context.Context, postID uuid.UUID) error {
	var entity entities.PostComment
	if err := r.db.WithContext(ctx).
		Model(&entity).
		Where("post_id = ?", postID).
		Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}
