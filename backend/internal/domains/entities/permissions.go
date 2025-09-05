package entities

import (
	"backend/internal/domains/commons"
)

type Permission struct {
	commons.Entity
	Name        string `gorm:"uniqueIndex;not null"`
	Description string `gorm:"null"`
	commons.Auditable

	RolePermission []RolePermission `gorm:"foreignKey:PermissionID"`
}
