package postgres

import (
	"backend/internal/persistents/abstractions"
	"context"

	"gorm.io/gorm"
)

type postUow struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewPostUow(db *gorm.DB) abstractions.PostUow {
	return &postUow{db: db}
}

func (u *postUow) Begin(ctx context.Context) (abstractions.PostRepoProvider, error) {
	tx := u.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	u.tx = tx
	return &postRepoProvider{tx: tx}, nil
}

func (u *postUow) Commit() error {
	if u.tx == nil {
		return nil
	}
	return u.tx.Commit().Error
}

func (u *postUow) Rollback() error {
	if u.tx == nil {
		return nil
	}
	return u.tx.Rollback().Error
}

func (u *postUow) GetDB() *gorm.DB {
	return u.db
}

type postRepoProvider struct {
	tx            *gorm.DB
	postRepo      abstractions.IPostRepository
	postImageRepo abstractions.IPostImageRepository
	postLikeRepo  abstractions.IPostLikeRepository
	commentRepo   abstractions.ICommentRepository
	eventRepo     abstractions.IEventRepository
}

func (r *postRepoProvider) PostRepository() abstractions.IPostRepository {
	if r.postRepo == nil {
		r.postRepo = NewPostRepo(r.tx)
	}
	return r.postRepo
}

func (r *postRepoProvider) PostImageRepository() abstractions.IPostImageRepository {
	if r.postImageRepo == nil {
		r.postImageRepo = NewPostImageRepo(r.tx)
	}
	return r.postImageRepo
}

func (r *postRepoProvider) PostLikeRepository() abstractions.IPostLikeRepository {
	if r.postLikeRepo == nil {
		r.postLikeRepo = NewPostLikeRepo(r.tx)
	}
	return r.postLikeRepo
}

func (r *postRepoProvider) CommentRepository() abstractions.ICommentRepository {
	if r.commentRepo == nil {
		r.commentRepo = NewCommentRepo(r.tx)
	}
	return r.commentRepo
}

// EventRepository implements abstractions.PostRepoProvider.
func (r *postRepoProvider) EventRepository() abstractions.IEventRepository {
	if r.eventRepo == nil {
		r.eventRepo = NewEventRepo(r.tx)
	}
	return r.eventRepo
}
