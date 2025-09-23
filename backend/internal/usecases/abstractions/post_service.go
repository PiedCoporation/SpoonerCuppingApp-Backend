package abstractions

import (
	"backend/internal/contracts/common"
	"backend/internal/contracts/post"
	"context"

	"github.com/google/uuid"
)

type IPostService interface {
	Create(ctx context.Context, userID uuid.UUID, req post.CreatePostReq) (uuid.UUID, error)
	GetAll(ctx context.Context, pageSize int, pageNumber int, searchTerm string) (*common.PageResult[post.PostViewRes], error)
	GetByID(ctx context.Context, id uuid.UUID) (*post.GetPostByIdRes, error)
	Update(ctx context.Context, userID, postID uuid.UUID, req post.UpdatePostReq) error
	TogglePostLike(ctx context.Context, userID, postID uuid.UUID) (*post.TogglePostLikeRes, error)
	Delete(ctx context.Context, userID, postID uuid.UUID) error
}
