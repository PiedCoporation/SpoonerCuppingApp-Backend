package mapper

import (
	postContract "backend/internal/contracts/post"
	"backend/internal/domains/entities"
)

// temp struct for gorm to scan
type PostWithCounts struct {
	entities.Post
	LikeCount    int64 `gorm:"->"`
	CommentCount int64 `gorm:"->"`
}

func MapPostToContractPostResponse(p *PostWithCounts) *postContract.PostViewRes {
	imageUrls := make([]string, len(p.Images))
	for j, img := range p.Images {
		imageUrls[j] = img.URL
	}

	return &postContract.PostViewRes{
		ID:           p.ID,
		Title:        p.Title,
		Content:      p.Content,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
		EventID:      p.EventID,
		User:         *MapUserToContractUserResponse(&p.User),
		ImageURLs:    imageUrls,
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
	}
}
