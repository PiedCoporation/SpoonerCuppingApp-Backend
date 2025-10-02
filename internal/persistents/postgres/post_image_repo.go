package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postImagePgRepo struct {
	*genericRepository[entities.PostImage]
}

func NewPostImageRepo(db *gorm.DB) abstractions.IPostImageRepository {
	return &postImagePgRepo{
		genericRepository: NewGenericRepository[entities.PostImage](db),
	}
}

// GetAllByPostID implements abstractions.IPostImageRepository.
func (p *postImagePgRepo) GetAllByPostID(ctx context.Context, postID uuid.UUID, preloads ...string) ([]entities.PostImage, error) {
	var entities []entities.PostImage
	db := p.db.WithContext(ctx)
	for _, p := range preloads {
		db = db.Preload(p)
	}
	if err := db.Where("post_id = ?", postID).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// DeleteByPostID implements abstractions.IPostImageRepository.
func (p *postImagePgRepo) DeleteByPostID(ctx context.Context, postID uuid.UUID) error {
	var entity entities.PostImage
	if err := p.db.WithContext(ctx).
		Model(&entity).
		Where("post_id = ?", postID).
		Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}

// DeleteByUrls implements abstractions.IPostImageRepository.
func (p *postImagePgRepo) DeleteByUrls(ctx context.Context, imageUrls []string) error {
	var entity entities.PostImage
	if err := p.db.WithContext(ctx).
		Model(&entity).
		Where("url IN ?", imageUrls).
		Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}
