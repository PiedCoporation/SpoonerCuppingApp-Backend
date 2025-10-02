package usecases

import (
	"backend/internal/constants/errorcode"
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
	postCmtRepo   persistentRepo.IPostCommentRepository
}

func NewPostService(
	postUow persistentRepo.PostUow,
	postRepo persistentRepo.IPostRepository,
	postImageRepo persistentRepo.IPostImageRepository,
	postLikeRepo persistentRepo.IPostLikeRepository,
	postCmtRepo persistentRepo.IPostCommentRepository,
) serviceAbstractions.IPostService {
	return &postService{
		postUow:       postUow,
		postRepo:      postRepo,
		postImageRepo: postImageRepo,
		postLikeRepo:  postLikeRepo,
		postCmtRepo:   postCmtRepo,
	}
}

// Create implements abstractions.IPostService.
func (s *postService) Create(ctx context.Context, userID uuid.UUID, req post.CreatePostReq) (uuid.UUID, error) {
	repoProvider, err := s.postUow.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer s.postUow.Rollback()

	postRepo := repoProvider.PostRepository()
	postImageRepo := repoProvider.PostImageRepository()
	eventRepo := repoProvider.EventRepository()

	// check event exists
	if req.EventID != nil {
		if _, err := eventRepo.GetByID(ctx, *req.EventID); err != nil {
			if errors.Is(err, errorcode.ErrNotFound) {
				return uuid.Nil, errorcode.ErrEventNotFound
			}
			return uuid.Nil, err
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
		return uuid.Nil, err
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
		return uuid.Nil, err
	}

	// commit the transaction
	if err := s.postUow.Commit(); err != nil {
		return uuid.Nil, err
	}

	return postEntity.ID, nil
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
		responses[i] = *mapper.MapPostToContractViewResponse(&p)
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
func (s *postService) GetByID(ctx context.Context, id uuid.UUID) (*post.GetPostByIdRes, error) {
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

	// get postcomment at first level
	rootComments, err := s.postCmtRepo.GetRootComments(ctx, id, true)
	if err != nil {
		return nil, err
	}

	commentResponses := make([]post.PostCommentViewRes, len(rootComments))
	for i, c := range rootComments {
		commentResponses[i] = *mapper.MapPostCommentToContractViewResponse(&c)
	}

	return &post.GetPostByIdRes{
		PostViewRes:  *mapper.MapPostToContractViewResponse(&postWithCounts),
		RootComments: commentResponses,
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

	fieldMap := make(map[string]any)

	if req.Title != nil && *req.Title != postEntity.Title {
		fieldMap["title"] = *req.Title
	}

	if req.Content != nil && *req.Content != postEntity.Content {
		fieldMap["content"] = *req.Content
	}

	if req.EventID != nil && (postEntity.EventID == nil || *req.EventID != *postEntity.EventID) {
		// check event exists
		if _, err := eventRepo.GetByID(ctx, *req.EventID); err != nil {
			if errors.Is(err, errorcode.ErrNotFound) {
				return errorcode.ErrEventNotFound
			}
			return err
		}
		fieldMap["event_id"] = *req.EventID
	}

	if err := postRepo.Update(ctx, postEntity.ID, fieldMap); err != nil {
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
	postCmtRepo := repoProvider.PostCommentRepository()

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
	// delete postcomment
	if err := postCmtRepo.DeleteByPostID(ctx, postID); err != nil {
		return err
	}

	return s.postUow.Commit()
}

// ====== Helper ======
func (s *postService) getPostWithCountsQuery(ctx context.Context, db *gorm.DB) *gorm.DB {
	// subquery for count likes and comments
	likeCount := db.Model(&entities.PostLike{}).
		Select("post_id, COUNT(*) as count").
		Where("is_deleted = ?", false).
		Group("post_id")
	commentCount := db.Model(&entities.PostComment{}).
		Select("post_id, COUNT(*) as count").
		Where("is_deleted = ?", false).
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
