package entities

import (
	"backend/internal/domains/commons"

	"github.com/google/uuid"
)

type PostImage struct {
	commons.Entity
	URL string `gorm:"not null"`
	commons.Auditable

	PostID uuid.UUID `gorm:"not null;index"`
	Post   Post      `gorm:"foreignKey:PostID"`
}
