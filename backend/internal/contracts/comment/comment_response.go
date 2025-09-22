package comment

import (
	"backend/internal/contracts/user"
	"time"

	"github.com/google/uuid"
)

type CommentViewRes struct {
	ID        uuid.UUID        `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Content   string           `json:"content" example:"Lorem ipsum dolor sit"`
	CreatedAt time.Time        `json:"created_at" example:"2024-01-15T10:00:00Z"`
	UpdatedAt time.Time        `json:"updated_at" example:"2024-01-15T10:00:00Z"`
	ParentID  *uuid.UUID       `json:"parent_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	User      user.UserViewRes `json:"user"`
}
