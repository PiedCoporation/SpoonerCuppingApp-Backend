package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postLikePgRepo struct {
	*genericRepository[entities.PostLike]
}

func NewPostLikeRepo(db *gorm.DB) abstractions.IPostLikeRepository {
	return &postLikePgRepo{
		genericRepository: NewGenericRepository[entities.PostLike](db),
	}
}

// DeleteByPostID implements abstractions.IPostLikeRepository.
func (p *postLikePgRepo) DeleteByPostID(ctx context.Context, postID uuid.UUID) error {
	var entity entities.PostLike
	if err := p.db.WithContext(ctx).
		Model(&entity).
		Where("post_id = ?", postID).
		Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}
