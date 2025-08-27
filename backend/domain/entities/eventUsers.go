package entities

import (
	"backend/domain/abstractions"

	"github.com/google/uuid"
)

type EventUser struct {
	abstractions.Entity
	IsAccepted bool `gorm:"not null"`
	IsInvited  bool `gorm:"not null"`
	abstractions.Auditable

	UserID  uuid.UUID `gorm:"not null;index"`
	User    User      `gorm:"foreignKey:UserID"`
	EventID uuid.UUID `gorm:"not null;index"`
	Event   Event     `gorm:"foreignKey:EventID"`
}
