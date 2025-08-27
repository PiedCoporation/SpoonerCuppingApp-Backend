package entities

import (
	"backend/domain/abstractions"
)

type Role struct {
	abstractions.Entity
	Name        string `gorm:"uniqueIndex;not null"`
	Description string `gorm:"null"`
	abstractions.Auditable

	RolePermission []RolePermission `gorm:"foreignKey:RoleID"`
}
