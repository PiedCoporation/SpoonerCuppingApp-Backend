package entities

import (
	"backend/internal/domains/commons"

	"github.com/google/uuid"
)

type PostComment struct {
	commons.Entity
	Content string `gorm:"not null"`
	commons.Auditable

	PostID uuid.UUID `gorm:"not null;index"`
	Post   Post      `gorm:"foreignKey:PostID"`

	ParentID *uuid.UUID   `gorm:"index"`
	Parent   *PostComment `gorm:"foreignKey:ParentID"`

	UserID uuid.UUID `gorm:"not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	Replies []PostComment `gorm:"foreignKey:ParentID"`
}
