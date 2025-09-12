package entities

import (
	"backend/internal/domains/commons"

	"github.com/google/uuid"
)

type EventAddress struct {
	commons.Entity
	Province  string `gorm:"not null"`
	District  string `gorm:"not null"`
	Longitude string `gorm:"not null"`
	Latitude  string `gorm:"not null"`
	Ward      string `gorm:"not null"`
	Street    string `gorm:"not null"`
	Phone     string `gorm:"not null"`
	commons.Auditable

	UserID uuid.UUID `gorm:"not null;index"`
	User   User      `gorm:"foreignKey:UserID"`
}
