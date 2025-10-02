package abstractions

import (
	"backend/internal/domains/entities"
	"context"

	"github.com/google/uuid"
)

type IPostImageRepository interface {
	IGenericRepository[entities.PostImage]
	GetAllByPostID(ctx context.Context, postID uuid.UUID, preloads ...string) ([]entities.PostImage, error)
	DeleteByPostID(ctx context.Context, postID uuid.UUID) error
	DeleteByUrls(ctx context.Context, imageUrls []string) error
}
