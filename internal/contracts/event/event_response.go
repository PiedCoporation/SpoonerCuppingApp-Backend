package event

import (
	"time"

	"backend/internal/constants/enums/eventregisterstatus"
	"backend/internal/constants/enums/processing"
	"backend/internal/constants/enums/roastinglever"

	"github.com/google/uuid"
)

// Event represents a cupping event
type Event struct{
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name string `json:"name" example:"Coffee Cupping Event"`
	DateOfEvent time.Time `json:"date_of_event" example:"2024-01-15T10:00:00Z"`
	StartTime time.Time `json:"start_time" example:"2024-01-15T10:00:00Z"`
	EndTime time.Time `json:"end_time" example:"2024-01-15T18:00:00Z"`
	Limit int `json:"limit" example:"50"`
	TotalCurrent int `json:"total_current" example:"20"`
	NumberSamples int `json:"number_samples" example:"5"`
	PhoneContact string `json:"phone_contact" example:"+1234567890"`
	EmailContact string `json:"email_contact" example:"contact@example.com"`
	IsPublic bool `json:"is_public" example:"true"`
	RegisterDate time.Time `json:"register_date" example:"2024-01-10T00:00:00Z"`
	RegisterStatus eventregisterstatus.RegisterStatusEnum `json:"register_status" example:"PENDING" enums:"PENDING,ACCEPTED,FULL"`
	EventAddress []EventAddress `json:"event_address"`
	HostBy HostBy `json:"host_by"`
}

// HostBy represents the host of an event
type HostBy struct {
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	FirstName string `json:"first_name" example:"John"`
	LastName string `json:"last_name" example:"Doe"`
	Email string `json:"email" example:"john.doe@example.com"`
	Phone string `json:"phone" example:"+1234567890"`
}

// EventAddress represents the address of an event
type EventAddress struct{
	Province string `json:"province" example:"Ho Chi Minh"`
	District string `json:"district" example:"District 1"`
	Longitude string `json:"longitude" example:"106.6297"`
	Latitude string `json:"latitude" example:"10.8231"`
	Ward string `json:"ward" example:"Ben Nghe Ward"`
	Street string `json:"street" example:"123 Main Street"`
	Phone string `json:"phone" example:"+1234567890"`
}

type GetEventByIDResponse struct {
	Event Event `json:"event"`
	Samples []EventSample `json:"samples"`
}

type EventSample struct {
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name string `json:"name" example:"Coffee Sample"`
	Rating *int `json:"rating" example:"5"`
	RoastingDate time.Time `json:"roasting_date" example:"2024-01-15T10:00:00Z"`
	RoastLevel roastinglever.RoastingLeverEnum `json:"roast_level" example:"Medium"`
	AltitudeGrow string `json:"altitude_grow" example:"1000m"`
	RoasteryName string `json:"roastery_name" example:"Roastery Name"`
	RoasteryAddress string `json:"roastery_address" example:"Roastery Address"`
	BreedName string `json:"breed_name" example:"Breed Name"`
	PreProcessing processing.ProcessingEnum `json:"pre_processing" example:"Washed"`
	GrowNation string `json:"grow_nation" example:"Vietnam"`
	GrowAddress string `json:"grow_address" example:"Grow Address"`
	Price *string `json:"price" example:"100000"`
}

// EventPageResult represents a paginated response for events
type EventPageResult struct {
	Data       []Event `json:"data" description:"Array of events"`
	Total      int     `json:"total" example:"150" description:"Total number of events"`
	Page       int     `json:"page" example:"1" description:"Current page number"`
	PageSize   int     `json:"page_size" example:"10" description:"Number of events per page"`
	TotalPages int     `json:"total_pages" example:"15" description:"Total number of pages"`
}