package entities

import (
	"backend/internal/domains/commons"

	"github.com/google/uuid"
)

type PostLike struct {
	commons.Entity
	commons.Auditable

	PostID uuid.UUID `gorm:"not null;index;uniqueIndex:idx_post_user"`
	Post   Post      `gorm:"foreignKey:PostID"`

	UserID uuid.UUID `gorm:"not null;index;uniqueIndex:idx_post_user"`
	User   User      `gorm:"foreignKey:UserID"`
}
