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

func MapPostToContractPostResponse(p *PostWithCounts) *postContract.PostResponse {
	imageUrls := make([]string, len(p.Images))
	for j, img := range p.Images {
		imageUrls[j] = img.URL
	}

	return &postContract.PostResponse{
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

func MapCommentToContractCommentResponse(c *entities.Comment) *postContract.CommentResponse {
	return &postContract.CommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		ParentID:  c.ParentID,
		User:      *MapUserToContractUserResponse(&c.User),
	}
}
