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
type LoginUserRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         UserRes `json:"user"`
}

type UserRes struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
	Role      string    `json:"role" example:"USER" enums:"USER,ADMIN"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Phone     string    `json:"phone" example:"1234567890"`
	// CircleStyle circlestyle.CircleStyleEnum `json:"circle_style" example:"DEFAULT" enums:"DEFAULT,SECOND"`
}