package entities

import (
	"backend/internal/domain/commons"
)

type Role struct {
	commons.Entity
	Name        string `gorm:"uniqueIndex;not null"`
	Description string `gorm:"null"`
	commons.Auditable

	Users          []User           `gorm:"foreignKey:RoleID"`
	RolePermission []RolePermission `gorm:"foreignKey:RoleID"`
}
