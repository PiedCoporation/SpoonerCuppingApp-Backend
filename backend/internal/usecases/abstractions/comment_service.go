package abstractions

import (
	"backend/internal/contracts/comment"
	"context"

	"github.com/google/uuid"
)

type ICommentService interface {
	Create(ctx context.Context, userID uuid.UUID, req comment.CreateCommentReq) error
	GetDirectChildren(ctx context.Context, parentID uuid.UUID, orderByCreatedAtDesc bool) ([]comment.CommentViewRes, error)
	GetRootCommentsByPostID(ctx context.Context, postID uuid.UUID, orderByCreatedAtDesc bool) ([]comment.CommentViewRes, error)
	Update(ctx context.Context, userID, commentID uuid.UUID, req comment.UpdateCommentReq) error
	Delete(ctx context.Context, userID, commentID uuid.UUID) error
}
