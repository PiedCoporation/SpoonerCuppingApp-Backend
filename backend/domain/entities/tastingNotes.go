package entities

import (
	"backend/domain/commons"

	"github.com/google/uuid"
)

type TastingNote struct {
	commons.Entity
	ParentName     string
	ChildName      string
	GrandChildName string
	commons.Auditable

	UserTastingID uuid.UUID   `gorm:"not null;index"`
	UserTasting   UserTasting `gorm:"foreignKey:UserTastingID"`
}
