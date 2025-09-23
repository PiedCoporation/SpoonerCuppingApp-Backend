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

	"github.com/google/uuid"
)

type postCommentService struct {
	postCommentRepo persistentRepo.IPostCommentRepository
	postRepo        persistentRepo.IPostRepository
}

func NewPostCommentService(
	postCommentRepo persistentRepo.IPostCommentRepository,
	postRepo persistentRepo.IPostRepository,
) serviceAbstractions.IPostCommentService {
	return &postCommentService{
		postCommentRepo: postCommentRepo,
		postRepo:        postRepo,
	}
}

// Create implements abstractions.ICommentService.
func (s *postCommentService) Create(ctx context.Context, userID, postID uuid.UUID, req post.CreatePostCommentReq) (uuid.UUID, error) {
	// get post
	if _, err := s.postRepo.GetByID(ctx, postID); err != nil {
		if errors.Is(err, errorcode.ErrNotFound) {
			return uuid.Nil, errorcode.ErrPostNotFound
		}
		return uuid.Nil, err
	}

	// check parent ID
	if req.ParentID != nil {
		if _, err := s.postCommentRepo.GetByID(ctx, *req.ParentID); err != nil {
			if errors.Is(err, errorcode.ErrNotFound) {
				return uuid.Nil, errorcode.ErrCommentNotFound
			}
			return uuid.Nil, err
		}
	}

	commentEntity := entities.PostComment{
		Entity:   commons.Entity{ID: uuid.New(), IsDeleted: false},
		Content:  req.Content,
		PostID:   postID,
		ParentID: req.ParentID,
		UserID:   userID,
	}

	if err := s.postCommentRepo.Create(ctx, &commentEntity); err != nil {
		return uuid.Nil, err
	}

	return commentEntity.ID, nil
}

// GetDirectChildren implements abstractions.ICommentService.
func (s *postCommentService) GetDirectChildren(
	ctx context.Context, parentID uuid.UUID, orderByCreatedAtDesc bool,
) ([]post.PostCommentViewRes, error) {
	// check parent exists
	if _, err := s.postCommentRepo.GetByID(ctx, parentID); err != nil {
		if errors.Is(err, errorcode.ErrNotFound) {
			return nil, errorcode.ErrCommentNotFound
		}
		return nil, err
	}

	comments, err := s.postCommentRepo.GetDirectChildren(ctx, parentID, orderByCreatedAtDesc)
	if err != nil {
		return nil, err
	}

	commentViews := make([]post.PostCommentViewRes, len(comments))
	for i, c := range comments {
		commentViews[i] = *mapper.MapPostCommentToContractViewResponse(&c)
	}

	return commentViews, nil
}

// GetRootCommentsByPostID implements abstractions.ICommentService.
func (s *postCommentService) GetRootCommentsByPostID(
	ctx context.Context, postID uuid.UUID, orderByCreatedAtDesc bool,
) ([]post.PostCommentViewRes, error) {
	// check parent exists
	if _, err := s.postRepo.GetByID(ctx, postID); err != nil {
		if errors.Is(err, errorcode.ErrNotFound) {
			return nil, errorcode.ErrPostNotFound
		}
		return nil, err
	}

	comments, err := s.postCommentRepo.GetRootComments(ctx, postID, orderByCreatedAtDesc)
	if err != nil {
		return nil, err
	}

	commentViews := make([]post.PostCommentViewRes, len(comments))
	for i, c := range comments {
		commentViews[i] = *mapper.MapPostCommentToContractViewResponse(&c)
	}

	return commentViews, nil
}

// Update implements abstractions.ICommentService.
func (s *postCommentService) Update(ctx context.Context, userID, commentID uuid.UUID, req post.UpdatePostCommentReq) error {
	if _, err := s.getOwnedComment(ctx, userID, commentID); err != nil {
		return err
	}

	if err := s.postCommentRepo.Update(ctx, commentID, map[string]any{
		"content": req.Content,
	}); err != nil {
		return err
	}

	return nil
}

// Delete implements abstractions.ICommentService.
func (s *postCommentService) Delete(ctx context.Context, userID, commentID uuid.UUID) error {
	if _, err := s.getOwnedComment(ctx, userID, commentID); err != nil {
		return err
	}

	if err := s.postCommentRepo.Update(ctx, commentID, map[string]any{
		"is_deleted": true,
	}); err != nil {
		return err
	}

	return nil
}

// helper
func (s *postCommentService) getOwnedComment(ctx context.Context, userID, commentID uuid.UUID) (*entities.PostComment, error) {
	commentEntity, err := s.postCommentRepo.GetByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, errorcode.ErrNotFound) {
			return nil, errorcode.ErrCommentNotFound
		}
		return nil, err
	}
	if commentEntity.UserID != userID {
		return nil, errorcode.ErrNotCommentOwner
	}
	return commentEntity, nil
}
