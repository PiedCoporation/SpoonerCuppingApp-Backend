package entities

import (
	"backend/internal/domains/commons"
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
	IsDefault bool
	commons.Auditable

	Events []Event `gorm:"foreignKey:EventAddressID"`
}
