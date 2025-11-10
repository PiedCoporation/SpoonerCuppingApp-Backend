package event

import (
	"github.com/google/uuid"
)

type GetAllUserEventResponse struct {
	UserEventID uuid.UUID `json:"user_event_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	FirstName string `json:"first_name" example:"John"`
	LastName string `json:"last_name" example:"Doe"`
	Email string `json:"email" example:"john.doe@example.com"`
	Phone string `json:"phone" example:"+1234567890"`
	IsAccepted bool `json:"is_accepted" example:"true"`
	IsInvited bool `json:"is_invited" example:"false"`
	IsHost bool `json:"is_host" example:"false"`
}