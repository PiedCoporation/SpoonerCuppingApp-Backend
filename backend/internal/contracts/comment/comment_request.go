package comment

import "github.com/google/uuid"

type CreateCommentReq struct {
	Content  string     `json:"comment" binding:"required"`
	PostID   uuid.UUID  `json:"post_id" binding:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
}

type UpdateCommentReq struct {
	Content string `json:"comment" binding:"required"`
}
