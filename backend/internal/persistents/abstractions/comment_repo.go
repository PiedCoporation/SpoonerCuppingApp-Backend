package abstractions

import (
	"backend/internal/domains/entities"
	"context"

	"github.com/google/uuid"
)

type ICommentRepository interface {
	GenericRepository[entities.Comment]
	GetDirectChildren(ctx context.Context, parentID uuid.UUID, orderByCreatedAtDesc bool) ([]entities.Comment, error)
	GetRootComments(ctx context.Context, postID uuid.UUID, orderByCreatedAtDesc bool) ([]entities.Comment, error)
	DeleteByPostID(ctx context.Context, postID uuid.UUID) error
}
