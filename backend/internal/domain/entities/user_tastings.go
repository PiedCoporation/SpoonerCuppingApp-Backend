package entities

import (
	"backend/internal/domain/commons"

	"github.com/google/uuid"
)

type UserTasting struct {
	commons.Entity
	rating int `gorm:"not null;check:rating_range,rating >= 1 AND rating <= 5"`
	commons.Auditable

	EventSampleID uuid.UUID   `gorm:"not null;index"`
	EventSample   EventSample `gorm:"foreignKey:EventSampleID"`
	EventUserID   uuid.UUID   `gorm:"not null;index"`
	EventUser     EventUser   `gorm:"foreignKey:EventUserID"`

	TastingNotes []TastingNote `gorm:"foreignKey:UserTastingID"`
}
