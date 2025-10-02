package usecases

import (
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/post"
	"backend/internal/domains/commons"
	"backend/internal/domains/entities"
	"backend/internal/mapper"
	persistentRepo "backend/internal/persistents/abstractions"
	serviceAbstractions "backend/internal/usecases/abstractions"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postLikeService struct {
	postRepo     persistentRepo.IPostRepository
	postLikeRepo persistentRepo.IPostLikeRepository
}

func NewPostLikeService(
	postRepo persistentRepo.IPostRepository,
	postLikeRepo persistentRepo.IPostLikeRepository,
) serviceAbstractions.IPostLikeService {
	return &postLikeService{
		postRepo:     postRepo,
		postLikeRepo: postLikeRepo,
	}
}

// GetPostLikeByPostID implements abstractions.IPostLikeService.
func (s *postLikeService) GetPostLikeByPostID(ctx context.Context, postID uuid.UUID) ([]post.PostLikeRes, error) {
	// check post exists
	_, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, errorcode.ErrNotFound) {
			return nil, errorcode.ErrPostNotFound
		}
		return nil, err
	}

	postLikes, err := s.postLikeRepo.FindByQuery(ctx, "post_id = ?", []any{postID}, true, "User")
	if err != nil {
		return nil, err
	}

	postLikesRes := make([]post.PostLikeRes, len(postLikes))
	for i, item := range postLikes {
		postLikesRes[i] = post.PostLikeRes{
			ID:   item.ID,
			User: *mapper.MapUserToContractUserResponse(&item.User),
		}
	}

	return postLikesRes, nil
}

// TogglePostLike implements abstractions.IPostService.
func (s *postLikeService) TogglePostLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (*post.TogglePostLikeRes, error) {
	// check post exists
	_, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, errorcode.ErrNotFound) {
			return nil, errorcode.ErrPostNotFound
		}
		return nil, err
	}

	query := fmt.Sprintf("user_id = '%s' AND post_id = '%s'", userID.String(), postID.String())
	likeEntity, err := s.postLikeRepo.GetSingle(ctx, query)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// create new if not exists
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newLikeEntity := &entities.PostLike{
			Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
			PostID: postID,
			UserID: userID,
		}
		if err := s.postLikeRepo.Create(ctx, newLikeEntity); err != nil {
			return nil, err
		}
		return &post.TogglePostLikeRes{
			ID:    newLikeEntity.ID,
			Liked: !newLikeEntity.IsDeleted,
		}, nil
	}

	// update is_deleted if exists
	if err := s.postLikeRepo.Update(ctx, likeEntity.ID, map[string]any{
		"is_deleted": !likeEntity.IsDeleted,
	}); err != nil {
		return nil, err
	}

	return &post.TogglePostLikeRes{
		ID:    likeEntity.ID,
		Liked: likeEntity.IsDeleted,
	}, nil
}
