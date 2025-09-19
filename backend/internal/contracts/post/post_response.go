package post

import (
	"time"

	"github.com/google/uuid"
)

// Post represents a post in newsfeed
type PostResponse struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title     string    `json:"title" example:"My title"`
	Content   string    `json:"content" example:"Lorem ipsum dolor sit amet consectetur adipiscing elit. Quisque faucibus ex sapien vitae pellentesque sem placerat. In id cursus mi pretium tellus duis convallis. Tempus leo eu aenean sed diam urna tempor. Pulvinar vivamus fringilla lacus nec metus bibendum egestas. Iaculis massa nisl malesuada lacinia integer nunc posuere. Ut hendrerit semper vel class aptent taciti sociosqu. Ad litora torquent per conubia nostra inceptos himenaeos."`
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-15T10:00:00Z"`

	EventID *uuid.UUID `json:"event_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	User    User       `json:"user"`

	ImageURLs    []string `json:"images,omitempty" example:"[\"https://example.com/img1.jpg\",\"https://example.com/img2.jpg\"]"`
	LikeCount    int64    `json:"like_count" example:"15"`
	CommentCount int64    `json:"comment_count" example:"4"`
}

// User represents the owner of post
type User struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Phone     string    `json:"phone" example:"+1234567890"`
}
