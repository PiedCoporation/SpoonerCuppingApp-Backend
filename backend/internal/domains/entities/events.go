package entities

import (
	"backend/internal/constants/enums/event"
	"backend/internal/domains/commons"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	commons.Entity
	Name          string    `gorm:"not null"`
	DateOfEvent   time.Time `gorm:"not null"`
	StartTime     time.Time `gorm:"not null"`
	EndTime       time.Time `gorm:"not null"`
	Limit         int       `gorm:"not null"`
	TotalCurrent  int       `gorm:"not null"`
	NumberSamples int       `gorm:"not null"`
	PhoneContact  string    `gorm:"not null"`
	EmailContact  string    `gorm:"not null"`
	InviteUrl     string    `gorm:"null"`
	QRImageUrl    string    `gorm:"null"`
	IsPublic      bool      `gorm:"not null"`
	RegisterDate  time.Time `gorm:"not null"`
	RegisterStatus event.RegisterStatusEnum `gorm:"not null"`
	commons.Auditable

	UserID         uuid.UUID    `gorm:"not null;index"`
	HostBy         User         `gorm:"foreignKey:UserID"`
	
	EventAddress   []EventAddress `gorm:"foreignKey:EventID"`
	EventUsers   []EventUser   `gorm:"foreignKey:EventID"`
	EventSamples []EventSample `gorm:"foreignKey:EventID"`
}
