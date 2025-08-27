package entities

import (
	"backend/domain/commons"

	"github.com/google/uuid"
)

type EventSample struct {
	commons.Entity
	Price  string `gorm:"not null"`
	Rating int    `gorm:"not null"`
	commons.Auditable

	UserSampleID uuid.UUID  `gorm:"not null;index"`
	UserSample   UserSample `gorm:"foreignKey:UserSampleID"`
	EventID      uuid.UUID  `gorm:"not null;index"`
	Event        Event      `gorm:"foreignKey:EventID"`

	UserTastings []UserTasting `gorm:"foreignKey:EventSampleID"`
}
