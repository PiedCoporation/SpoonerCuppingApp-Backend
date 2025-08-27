package entities

import (
	"backend/domain/commons"

	"github.com/google/uuid"
)

type RolePermission struct {
	commons.Entity
	commons.Auditable

	PermissionID uuid.UUID  `gorm:"not null;index"`
	Permission   Permission `gorm:"foreignKey:PermissionID"`
	RoleID       uuid.UUID  `gorm:"not null;index"`
	Role         Role       `gorm:"foreignKey:RoleID"`
}
