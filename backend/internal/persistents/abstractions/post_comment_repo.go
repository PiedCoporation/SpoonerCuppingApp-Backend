package abstractions

import (
	"backend/internal/domains/entities"
	"context"

	"github.com/google/uuid"
)

type IPostCommentRepository interface {
	GenericRepository[entities.PostComment]
	GetDirectChildren(ctx context.Context, parentID uuid.UUID, orderByCreatedAtDesc bool) ([]entities.PostComment, error)
	GetRootComments(ctx context.Context, postID uuid.UUID, orderByCreatedAtDesc bool) ([]entities.PostComment, error)
	DeleteByPostID(ctx context.Context, postID uuid.UUID) error
}
