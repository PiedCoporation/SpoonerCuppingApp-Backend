package entities

import (
	"backend/internal/domains/commons"

	"github.com/google/uuid"
)

type Post struct {
	commons.Entity
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	commons.Auditable

	EventID *uuid.UUID `gorm:"null;index"`
	Event   Event      `gorm:"foreignKey:EventID"`

	UserID uuid.UUID `gorm:"not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	Images   []PostImage   `gorm:"foreignKey:PostID"`
	Likes    []PostLike    `gorm:"foreignKey:PostID"`
	Comments []PostComment `gorm:"foreignKey:PostID"`
}
