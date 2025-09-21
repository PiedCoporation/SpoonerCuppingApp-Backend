package abstractions

import (
	"backend/internal/domains/entities"
	"context"

	"github.com/google/uuid"
)

type ICommentRepository interface {
	GenericRepository[entities.Comment]
	DeleteByPostID(ctx context.Context, postID uuid.UUID) error
}
