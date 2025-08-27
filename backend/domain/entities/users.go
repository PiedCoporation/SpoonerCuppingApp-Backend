package entities

import (
	"backend/domain/abstractions"
)

type User struct {
	abstractions.Entity
	FirstName  string `gorm:"null"`
	LastName   string `gorm:"null"`
	Username   string `gorm:"null"`
	Email      string `gorm:"uniqueIndex;not null"`
	Phone      string `gorm:"null"`
	IsVerified bool   `gorm:"default:false"`
	abstractions.Auditable

	Events     []Event     `gorm:"foreignKey:UserID"`
	EventUsers []EventUser `gorm:"foreignKey:UserID"`
}
