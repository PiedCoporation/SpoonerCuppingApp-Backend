package abstractions

import (
	"backend/internal/contracts/post"
	"context"

	"github.com/google/uuid"
)

type IPostCommentService interface {
	Create(ctx context.Context, userID, postID uuid.UUID, req post.CreatePostCommentReq) (uuid.UUID, error)
	GetDirectChildren(ctx context.Context, parentID uuid.UUID, orderByCreatedAtDesc bool) ([]post.PostCommentViewRes, error)
	GetRootCommentsByPostID(ctx context.Context, postID uuid.UUID, orderByCreatedAtDesc bool) ([]post.PostCommentViewRes, error)
	Update(ctx context.Context, userID, commentID uuid.UUID, req post.UpdatePostCommentReq) error
	Delete(ctx context.Context, userID, commentID uuid.UUID) error
}
