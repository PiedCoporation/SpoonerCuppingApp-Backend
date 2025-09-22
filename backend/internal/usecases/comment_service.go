package usecases

import (
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/comment"
	"backend/internal/domains/commons"
	"backend/internal/domains/entities"
	"backend/internal/mapper"
	persistentRepo "backend/internal/persistents/abstractions"
	serviceAbstractions "backend/internal/usecases/abstractions"
	"context"
	"errors"

	"github.com/google/uuid"
)

type commentService struct {
	commentRepo persistentRepo.ICommentRepository
	postRepo    persistentRepo.IPostRepository
}

func NewCommentService(
	commentRepo persistentRepo.ICommentRepository,
	postRepo persistentRepo.IPostRepository,
) serviceAbstractions.ICommentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

// Create implements abstractions.ICommentService.
func (s *commentService) Create(ctx context.Context, userID uuid.UUID, req comment.CreateCommentReq) error {
	// get post
	if _, err := s.postRepo.GetByID(ctx, req.PostID); err != nil {
		if errors.Is(err, errorcode.ErrNotFound) {
			return errorcode.ErrPostNotFound
		}
		return err
	}

	// check parent ID
	if req.ParentID != nil {
		if _, err := s.commentRepo.GetByID(ctx, *req.ParentID); err != nil {
			if errors.Is(err, errorcode.ErrNotFound) {
				return errorcode.ErrCommentNotFound
			}
			return err
		}
	}

	commentEntity := entities.Comment{
		Entity:   commons.Entity{ID: uuid.New(), IsDeleted: false},
		Content:  req.Content,
		PostID:   req.PostID,
		ParentID: req.ParentID,
		UserID:   userID,
	}

	if err := s.commentRepo.Create(ctx, &commentEntity); err != nil {
		return err
	}

	return nil
}

// GetDirectChildren implements abstractions.ICommentService.
func (s *commentService) GetDirectChildren(
	ctx context.Context, parentID uuid.UUID, orderByCreatedAtDesc bool,
) ([]comment.CommentViewRes, error) {
	comments, err := s.commentRepo.GetDirectChildren(ctx, parentID, orderByCreatedAtDesc)
	if err != nil {
		return nil, err
	}

	commentViews := make([]comment.CommentViewRes, len(comments))
	for i, c := range comments {
		commentViews[i] = *mapper.MapCommentToContractCommentResponse(&c)
	}

	return commentViews, nil
}

// GetRootCommentsByPostID implements abstractions.ICommentService.
func (s *commentService) GetRootCommentsByPostID(
	ctx context.Context, postID uuid.UUID, orderByCreatedAtDesc bool,
) ([]comment.CommentViewRes, error) {
	comments, err := s.commentRepo.GetRootComments(ctx, postID, orderByCreatedAtDesc)
	if err != nil {
		return nil, err
	}

	commentViews := make([]comment.CommentViewRes, len(comments))
	for i, c := range comments {
		commentViews[i] = *mapper.MapCommentToContractCommentResponse(&c)
	}

	return commentViews, nil
}

// Update implements abstractions.ICommentService.
func (s *commentService) Update(ctx context.Context, userID, commentID uuid.UUID, req comment.UpdateCommentReq) error {
	if _, err := s.getOwnedComment(ctx, userID, commentID); err != nil {
		return err
	}

	if err := s.commentRepo.Update(ctx, commentID, map[string]any{
		"content": req.Content,
	}); err != nil {
		return err
	}

	return nil
}

// Delete implements abstractions.ICommentService.
func (s *commentService) Delete(ctx context.Context, userID, commentID uuid.UUID) error {
	if _, err := s.getOwnedComment(ctx, userID, commentID); err != nil {
		return err
	}

	if err := s.commentRepo.Update(ctx, commentID, map[string]any{
		"is_deleted": true,
	}); err != nil {
		return err
	}

	return nil
}

// helper
func (s *commentService) getOwnedComment(ctx context.Context, userID, commentID uuid.UUID) (*entities.Comment, error) {
	commentEntity, err := s.commentRepo.GetByID(ctx, commentID)
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
