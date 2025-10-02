package abstractions

import (
	"context"

	"gorm.io/gorm"
)

type PostUow interface {
	Begin(ctx context.Context) (PostRepoProvider, error)
	Commit() error
	Rollback() error
	GetDB() *gorm.DB
}

type PostRepoProvider interface {
	PostRepository() IPostRepository
	PostImageRepository() IPostImageRepository
	PostLikeRepository() IPostLikeRepository
	PostCommentRepository() IPostCommentRepository
	EventRepository() IEventRepository
}
