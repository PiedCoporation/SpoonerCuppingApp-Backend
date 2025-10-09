package entities

import (
	"backend/internal/domains/commons"

	"github.com/google/uuid"
)

type EventUser struct {
	commons.Entity
	IsAccepted bool `gorm:"not null"`
	IsInvited  bool `gorm:"not null"`
	commons.Auditable

	UserID  uuid.UUID `gorm:"not null;index"`
	User    User      `gorm:"foreignKey:UserID"`
	EventID uuid.UUID `gorm:"not null;index"`
	Event   Event     `gorm:"foreignKey:EventID"`
	IsHost   bool      `gorm:"not null"`

	UserTastings []UserTasting `gorm:"foreignKey:EventUserID"`
}
