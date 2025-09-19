package usecases

import (
	"backend/internal/contracts/common"
	"backend/internal/contracts/post"
	"backend/internal/domains/entities"
	"backend/internal/mapper"
	persistentRepo "backend/internal/persistents/abstractions"
	serviceAbstractions "backend/internal/usecases/abstractions"
	"context"
	"math"
)

type postService struct {
	postUow  persistentRepo.EventUOW
	postRepo persistentRepo.IPostRepository
}

func NewPostService(
	postUow persistentRepo.EventUOW,
	postRepo persistentRepo.IPostRepository,
) serviceAbstractions.IPostService {
	return &postService{
		postUow:  postUow,
		postRepo: postRepo,
	}
}

// GetAll implements abstractions.IPostService.
func (s *postService) GetAll(ctx context.Context,
	pageSize int, pageNumber int, searchTerm string,
) (*common.PageResult[post.PostResponse], error) {
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

	q := db.WithContext(ctx).Model(&entities.Post{})

	if searchTerm != "" {
		q = q.Where("title ILIKE ?", "%"+searchTerm+"%")
	}
	q = q.Where("posts.is_deleted = ?", false)

	// subquery for count likes and comments
	likeCount := db.Model(&entities.PostLike{}).
		Select("post_id, COUNT(*) as count").
		Group("post_id")
	commentCount := db.Model(&entities.Comment{}).
		Select("post_id, COUNT(*) as count").
		Group("post_id")

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	var postWithCounts []mapper.PostWithCounts
	err := q.
		Select(`posts.*,
			COALESCE(lc.count, 0) as like_count,
			COALESCE(cc.count, 0) as comment_count`).
		Joins("LEFT JOIN (?) AS lc on lc.post_id = posts.id", likeCount).
		Joins("LEFT JOIN (?) AS cc on cc.post_id = posts.id", commentCount).
		Order("posts.created_at DESC").
		Offset(offset).Limit(pageSize).
		Preload("User").Preload("Images").
		Scan(&postWithCounts).Error
	if err != nil {
		return nil, err
	}

	responses := make([]post.PostResponse, len(postWithCounts))
	for i, p := range postWithCounts {
		responses[i] = *mapper.MapPostToContractGetAllPostResponse(p)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &common.PageResult[post.PostResponse]{
		Data:       responses,
		Total:      int(total),
		Page:       pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
