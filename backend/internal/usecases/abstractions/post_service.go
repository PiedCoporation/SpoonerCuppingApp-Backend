package abstractions

import (
	"backend/internal/contracts/common"
	"backend/internal/contracts/post"
	"context"

	"github.com/google/uuid"
)

type IPostService interface {
	Create(ctx context.Context, userID uuid.UUID, req post.CreatePostReq) error
	GetAll(ctx context.Context, pageSize int, pageNumber int, searchTerm string) (*common.PageResult[post.PostResponse], error)
	GetByID(ctx context.Context, id uuid.UUID) (*post.GetPostByIdResponse, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req post.UpdatePostReq) error
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}
