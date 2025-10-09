package entities

import (
	"backend/internal/domains/commons"

	"github.com/google/uuid"
)

type EventSample struct {
	commons.Entity
	Price  *string `gorm:"null"`
	Rating *int    `gorm:"null"`
	commons.Auditable

	UserSampleID uuid.UUID  `gorm:"not null;index"`
	UserSample   UserSample `gorm:"foreignKey:UserSampleID"`
	EventID      uuid.UUID  `gorm:"not null;index"`
	Event        Event      `gorm:"foreignKey:EventID"`

	UserTastings []UserTasting `gorm:"foreignKey:EventSampleID"`
}
