package user

import (
	"backend/internal/constants/enums/circlestyle"

	"github.com/google/uuid"
)

type UserViewRes struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
}

type UserSettingRes struct {
	ID          uuid.UUID                   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CircleStyle circlestyle.CircleStyleEnum `json:"circle_style" example:"DEFAULT" enums:"DEFAULT,SECOND"`
}
