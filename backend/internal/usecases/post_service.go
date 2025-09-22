package usecases

import (
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/comment"
	"backend/internal/contracts/common"
	"backend/internal/contracts/post"
	"backend/internal/domains/commons"
	"backend/internal/domains/entities"
	"backend/internal/mapper"
	persistentRepo "backend/internal/persistents/abstractions"
	serviceAbstractions "backend/internal/usecases/abstractions"
	"backend/pkg/utils/arrayutils"
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postService struct {
	postUow       persistentRepo.PostUow
	postRepo      persistentRepo.IPostRepository
	postImageRepo persistentRepo.IPostImageRepository
	postLikeRepo  persistentRepo.IPostLikeRepository
	commentRepo   persistentRepo.ICommentRepository
}

func NewPostService(
	postUow persistentRepo.PostUow,
	postRepo persistentRepo.IPostRepository,
	postImageRepo persistentRepo.IPostImageRepository,
	postLikeRepo persistentRepo.IPostLikeRepository,
	commentRepo persistentRepo.ICommentRepository,
) serviceAbstractions.IPostService {
	return &postService{
		postUow:       postUow,
		postRepo:      postRepo,
		postImageRepo: postImageRepo,
		postLikeRepo:  postLikeRepo,
		commentRepo:   commentRepo,
	}
}

// Create implements abstractions.IPostService.
func (s *postService) Create(ctx context.Context, userID uuid.UUID, req post.CreatePostReq) error {
	repoProvider, err := s.postUow.Begin(ctx)
	if err != nil {
		return err
	}
	defer s.postUow.Rollback()

	postRepo := repoProvider.PostRepository()
	postImageRepo := repoProvider.PostImageRepository()
	eventRepo := repoProvider.EventRepository()

	// check event exists
	if req.EventID != nil {
		if _, err := eventRepo.GetByID(ctx, *req.EventID); err != nil {
			if errors.Is(err, errorcode.ErrNotFound) {
				return errorcode.ErrEventNotFound
			}
			return err
		}
	}

	// create post entity
	postEntity := entities.Post{
		Entity:  commons.Entity{ID: uuid.New(), IsDeleted: false},
		Title:   req.Title,
		Content: req.Content,
		EventID: req.EventID,
		UserID:  userID,
	}

	// insert post to db
	if err := postRepo.Create(ctx, &postEntity); err != nil {
		return err
	}

	// make images slice
	images := make([]entities.PostImage, len(req.ImageUrls))
	for i, url := range req.ImageUrls {
		images[i] = entities.PostImage{
			Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
			URL:    url,
			PostID: postEntity.ID,
		}
	}

	// insert into db
	if err := postImageRepo.CreateRange(ctx, images); err != nil {
		return err
	}

	// commit the transaction
	if err := s.postUow.Commit(); err != nil {
		return err
	}

	return nil
}

// GetAll implements abstractions.IPostService.
func (s *postService) GetAll(ctx context.Context, pageSize int, pageNumber int, searchTerm string,
) (*common.PageResult[post.PostViewRes], error) {
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (pageNumber - 1) * pageSize

	db := s.postUow.GetDB()
	q := s.getPostWithCountsQuery(ctx, db)

	if searchTerm != "" {
		q = q.Where("title ILIKE ?", "%"+searchTerm+"%")
	}
	q = q.Where("posts.is_deleted = ?", false)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	var postWithCounts []mapper.PostWithCounts
	err := q.
		Order("posts.created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&postWithCounts).Error
	if err != nil {
		return nil, err
	}

	responses := make([]post.PostViewRes, len(postWithCounts))
	for i, p := range postWithCounts {
		responses[i] = *mapper.MapPostToContractPostResponse(&p)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &common.PageResult[post.PostViewRes]{
		Data:       responses,
		Total:      int(total),
		Page:       pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByID implements abstractions.IPostService.
func (s *postService) GetByID(ctx context.Context, id uuid.UUID) (*post.GetPostByIdResponse, error) {
	db := s.postUow.GetDB()

	var postWithCounts mapper.PostWithCounts
	err := s.getPostWithCountsQuery(ctx, db).
		Where("posts.id = ? AND posts.is_deleted = ?", id, false).
		First(&postWithCounts).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorcode.ErrPostNotFound
		}
		return nil, err
	}

	// get comment at first level
	var parentComments []entities.Comment
	err = db.WithContext(ctx).Model(&entities.Comment{}).
		Where("post_id = ? AND parent_id IS NULL", id).
		Preload("User").
		Order("created_at DESC").
		Find(&parentComments).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	commentResponses := make([]comment.CommentViewRes, len(parentComments))
	for i, c := range parentComments {
		commentResponses[i] = *mapper.MapCommentToContractCommentResponse(&c)
	}

	return &post.GetPostByIdResponse{
		PostViewRes:    *mapper.MapPostToContractPostResponse(&postWithCounts),
		ParentComments: commentResponses,
	}, nil
}

// Update implements abstractions.IPostService.
func (s *postService) Update(ctx context.Context, userID, postID uuid.UUID, req post.UpdatePostReq) error {
	repoProvider, err := s.postUow.Begin(ctx)
	if err != nil {
		return err
	}
	defer s.postUow.Rollback()

	postRepo := repoProvider.PostRepository()
	postImageRepo := repoProvider.PostImageRepository()
	eventRepo := repoProvider.EventRepository()

	// get post from db
	postEntity, err := s.getOwnedPost(ctx, userID, userID)
	if err != nil {
		return err
	}

	if req.Title != nil {
		postEntity.Title = *req.Title
	}
	if req.Content != nil {
		postEntity.Content = *req.Content
	}
	if req.EventID != nil {
		if _, err := eventRepo.GetByID(ctx, *req.EventID); err != nil {
			if errors.Is(err, errorcode.ErrNotFound) {
				return errorcode.ErrEventNotFound
			}
			return err
		}
		postEntity.EventID = req.EventID
	}

	if err := postRepo.Update(ctx, postEntity.ID, map[string]any{
		"title":    postEntity.Title,
		"content":  postEntity.Content,
		"event_id": postEntity.EventID,
	}); err != nil {
		return err
	}

	if req.ImageUrls != nil {
		oldImages, err := postImageRepo.GetAllByPostID(ctx, postEntity.ID)
		if err != nil {
			return err
		}
		oldUrls := make([]string, len(oldImages))
		for i, img := range oldImages {
			oldUrls[i] = img.URL
		}

		newUrls := *req.ImageUrls

		// remove image
		removed := arrayutils.Difference(oldUrls, newUrls)
		if err := postImageRepo.DeleteByUrls(ctx, removed); err != nil {
			return err
		}

		// add new image
		added := arrayutils.Difference(newUrls, oldUrls)
		// make images slice
		images := make([]entities.PostImage, len(added))
		for i, url := range added {
			images[i] = entities.PostImage{
				Entity: commons.Entity{ID: uuid.New(), IsDeleted: false},
				URL:    url,
				PostID: postEntity.ID,
			}
		}
		if err := postImageRepo.CreateRange(ctx, images); err != nil {
			return err
		}
	}

	return s.postUow.Commit()
}

// Delete implements abstractions.IPostService.
func (s *postService) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	repoProvider, err := s.postUow.Begin(ctx)
	if err != nil {
		return err
	}
	defer s.postUow.Rollback()

	postRepo := repoProvider.PostRepository()
	postImageRepo := repoProvider.PostImageRepository()
	postLikeRepo := repoProvider.PostLikeRepository()
	commentRepo := repoProvider.CommentRepository()

	// get post from db
	if _, err := s.getOwnedPost(ctx, userID, userID); err != nil {
		return err
	}

	// delete post
	if err := postRepo.SoftDelete(ctx, postID); err != nil {
		return err
	}

	// delete post image
	if err := postImageRepo.DeleteByPostID(ctx, postID); err != nil {
		return err
	}
	// delete post like
	if err := postLikeRepo.DeleteByPostID(ctx, postID); err != nil {
		return err
	}
	// delete comment
	if err := commentRepo.DeleteByPostID(ctx, postID); err != nil {
		return err
	}

	return s.postUow.Commit()
}

// ====== Helper ======
func (s *postService) getPostWithCountsQuery(ctx context.Context, db *gorm.DB) *gorm.DB {
	// subquery for count likes and comments
	likeCount := db.Model(&entities.PostLike{}).
		Select("post_id, COUNT(*) as count").
		Group("post_id")
	commentCount := db.Model(&entities.Comment{}).
		Select("post_id, COUNT(*) as count").
		Group("post_id")

	q := db.WithContext(ctx).Model(&entities.Post{}).
		Select(`posts.*,
			COALESCE(lc.count, 0) as like_count,
			COALESCE(cc.count, 0) as comment_count`).
		Joins("LEFT JOIN (?) AS lc on lc.post_id = posts.id", likeCount).
		Joins("LEFT JOIN (?) AS cc on cc.post_id = posts.id", commentCount).
		Preload("User").Preload("Images")

	return q
}

func (s *postService) getOwnedPost(ctx context.Context, userID, postID uuid.UUID) (*entities.Post, error) {
	postEntity, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, errorcode.ErrNotFound) {
			return nil, errorcode.ErrPostNotFound
		}
		return nil, err
	}
	if postEntity.UserID != userID {
		return nil, errorcode.ErrNotPostOwner
	}
	return postEntity, nil
}
