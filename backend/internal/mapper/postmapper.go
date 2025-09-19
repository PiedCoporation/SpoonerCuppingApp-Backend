package mapper

import (
	postContract "backend/internal/contracts/post"
	"backend/internal/domains/entities"
)

// temp struct for gorm to scan
type PostWithCounts struct {
	entities.Post
	LikeCount    int64
	CommentCount int64
}

func MapPostToContractGetAllPostResponse(p PostWithCounts) *postContract.PostResponse {
	imageUrls := make([]string, len(p.Images))
	for j, img := range p.Images {
		imageUrls[j] = img.URL
	}

	return &postContract.PostResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		EventID:   p.EventID,
		User: postContract.User{
			ID:        p.User.ID,
			FirstName: p.User.FirstName,
			LastName:  p.User.LastName,
			Email:     p.User.Email,
			Phone:     p.User.Phone,
		},
		ImageURLs:    imageUrls,
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
	}
}
