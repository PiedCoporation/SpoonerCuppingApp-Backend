package entities

import (
	"backend/domain/commons"

	"github.com/google/uuid"
)

type User struct {
	commons.Entity
	FirstName  string `gorm:"null"`
	LastName   string `gorm:"null"`
	Username   string `gorm:"null"`
	Email      string `gorm:"uniqueIndex;not null"`
	Phone      string `gorm:"null"`
	IsVerified bool   `gorm:"default:false"`
	commons.Auditable

	RoleID uuid.UUID `gorm:"not null;index"`
	Role   Role      `gorm:"foreignKey:RoleID"`

	Events      []Event      `gorm:"foreignKey:UserID"`
	EventUsers  []EventUser  `gorm:"foreignKey:UserID"`
	UserSamples []UserSample `gorm:"foreignKey:UserID"`
}
