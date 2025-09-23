package post

import "github.com/google/uuid"

type CreatePostReq struct {
	Title     string     `json:"title" binding:"required"`
	Content   string     `json:"content" binding:"required"`
	EventID   *uuid.UUID `json:"event_id"`
	ImageUrls []string   `json:"image_urls"`
}

type UpdatePostReq struct {
	Title     *string    `json:"title"`
	Content   *string    `json:"content"`
	EventID   *uuid.UUID `json:"event_id"`
	ImageUrls *[]string  `json:"image_urls"`
}

type CreatePostCommentReq struct {
	Content  string     `json:"content" binding:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
}

type UpdatePostCommentReq struct {
	Content string `json:"content" binding:"required"`
}
