package entities

import (
	"backend/internal/constants/enum/processing"
	"backend/internal/constants/enum/roastinglever"
	"backend/internal/domains/commons"
	"time"

	"github.com/google/uuid"
)

type UserSample struct {
	commons.Entity
	Name            string                          `gorm:"not null"`
	RoastingDate    time.Time                       `gorm:"not null"`
	RoastLevel      roastinglever.RoastingLeverEnum `gorm:"not null"`
	AltitudeGrow    string                          `gorm:"not null"`
	RoasteryName    string                          `gorm:"not null"`
	RoasteryAddress string                          `gorm:"not null"`
	BreedName       string                          `gorm:"not null"`
	PreProcessing   processing.ProcessingEnum       `gorm:"not null"`
	GrowNation      string                          `gorm:"not null"`
	GrowAddress     string                          `gorm:"not null"`
	Price           float64                         `gorm:"not null"`
	commons.Auditable

	UserID uuid.UUID `gorm:"not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	EventSamples []EventSample `gorm:"foreignKey:UserSampleID"`
}
