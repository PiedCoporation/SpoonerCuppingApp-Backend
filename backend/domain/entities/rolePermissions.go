package entities

import (
	"backend/domain/abstractions"

	"github.com/google/uuid"
)

type RolePermission struct {
	abstractions.Entity
	abstractions.Auditable

	PermissionID uuid.UUID  `gorm:"not null;index"`
	Permission   Permission `gorm:"foreignKey:PermissionID"`
	RoleID       uuid.UUID  `gorm:"not null;index"`
	Role         Role       `gorm:"foreignKey:RoleID"`
}
