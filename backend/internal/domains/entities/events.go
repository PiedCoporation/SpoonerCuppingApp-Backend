package entities

import (
	"backend/internal/constants/enums/event"
	eventContract "backend/internal/contracts/event"
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

// ToContract converts the Event entity to the contract response type
func (e *Event) ToContract() eventContract.Event {
	var addresses []eventContract.EventAddress
	for _, addr := range e.EventAddress {
		addresses = append(addresses, eventContract.EventAddress{
			Province:  addr.Province,
			District:  addr.District,
			Longitude: addr.Longitude,
			Latitude:  addr.Latitude,
			Ward:      addr.Ward,
			Street:    addr.Street,
			Phone:     addr.Phone,
		})
	}

	return eventContract.Event{
		ID:             e.ID,
		Name:           e.Name,
		DateOfEvent:    e.DateOfEvent,
		StartTime:      e.StartTime,
		EndTime:        e.EndTime,
		Limit:          e.Limit,
		NumberSamples:  e.NumberSamples,
		PhoneContact:   e.PhoneContact,
		EmailContact:   e.EmailContact,
		IsPublic:       e.IsPublic,
		RegisterDate:   e.RegisterDate,
		RegisterStatus: eventContract.RegisterStatusEnum(e.RegisterStatus),
		EventAddress:   addresses,
	}
}
