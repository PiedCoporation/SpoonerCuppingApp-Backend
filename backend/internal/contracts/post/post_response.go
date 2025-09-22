package post

import (
	"backend/internal/contracts/comment"
	"backend/internal/contracts/user"
	"time"

	"github.com/google/uuid"
)

// Post represents a post in newsfeed
type PostViewRes struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title     string    `json:"title" example:"My title"`
	Content   string    `json:"content" example:"Lorem ipsum dolor sit"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-15T10:00:00Z"`

	EventID *uuid.UUID       `json:"event_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	User    user.UserViewRes `json:"user"`

	ImageURLs    []string `json:"images" example:"[\"https://example.com/img1.jpg\"]"`
	LikeCount    int64    `json:"like_count" example:"15"`
	CommentCount int64    `json:"comment_count" example:"4"`
}

type GetPostByIdResponse struct {
	PostViewRes
	ParentComments []comment.CommentViewRes `json:"parent_comments"`
}

// PostPageResult represents a paginated response for posts
type PostPageResult struct {
	Data       []PostViewRes `json:"data" description:"Array of posts"`
	Total      int           `json:"total" example:"150" description:"Total number of posts"`
	Page       int           `json:"page" example:"1" description:"Current page number"`
	PageSize   int           `json:"page_size" example:"10" description:"Number of posts per page"`
	TotalPages int           `json:"total_pages" example:"15" description:"Total number of pages"`
}
