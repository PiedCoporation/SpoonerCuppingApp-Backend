package abstractions

import (
	"backend/internal/domains/entities"
	"context"

	"github.com/google/uuid"
)

type IPostLikeRepository interface {
	IGenericRepository[entities.PostLike]
	DeleteByPostID(ctx context.Context, postID uuid.UUID) error
}
