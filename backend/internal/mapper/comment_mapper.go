package mapper

import (
	commentContract "backend/internal/contracts/comment"
	"backend/internal/domains/entities"
)

func MapCommentToContractCommentResponse(c *entities.Comment) *commentContract.CommentViewRes {
	return &commentContract.CommentViewRes{
		ID:        c.ID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		ParentID:  c.ParentID,
		User:      *MapUserToContractUserResponse(&c.User),
	}
}
