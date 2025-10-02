package abstractions

import (
	"backend/internal/contracts/common"
	"backend/internal/contracts/post"
	"context"

	"github.com/google/uuid"
)

type (
	IPostService interface {
		Create(ctx context.Context, userID uuid.UUID, req post.CreatePostReq) (uuid.UUID, error)
		GetAll(ctx context.Context, pageSize int, pageNumber int, searchTerm string) (*common.PageResult[post.PostViewRes], error)
		GetByID(ctx context.Context, id uuid.UUID) (*post.GetPostByIdRes, error)
		Update(ctx context.Context, userID, postID uuid.UUID, req post.UpdatePostReq) error
		Delete(ctx context.Context, userID, postID uuid.UUID) error
	}

	IPostLikeService interface {
		GetPostLikeByPostID(ctx context.Context, postID uuid.UUID) ([]post.PostLikeRes, error)
		TogglePostLike(ctx context.Context, userID, postID uuid.UUID) (*post.TogglePostLikeRes, error)
	}

	IPostCommentService interface {
		Create(ctx context.Context, userID, postID uuid.UUID, req post.CreatePostCommentReq) (uuid.UUID, error)
		GetDirectChildren(ctx context.Context, parentID uuid.UUID, orderByCreatedAtDesc bool) ([]post.PostCommentViewRes, error)
		GetRootCommentsByPostID(ctx context.Context, postID uuid.UUID, orderByCreatedAtDesc bool) ([]post.PostCommentViewRes, error)
		Update(ctx context.Context, userID, commentID uuid.UUID, req post.UpdatePostCommentReq) error
		Delete(ctx context.Context, userID, commentID uuid.UUID) error
	}
)
