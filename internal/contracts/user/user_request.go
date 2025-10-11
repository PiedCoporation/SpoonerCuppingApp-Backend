package user

import "backend/internal/constants/enums/circlestyle"

type RegisterUserReq struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Phone           string `json:"phone" binding:"required"`
	Password        string `json:"password" binding:"required,min=8,max=30"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type ResendEmailReq struct {
	Email string `json:"email" binding:"required,email"`
}

type LoginUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=30"`
}

type LogoutUserReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordReq struct {
	Email string `json:"email" binding:"required,email"`
}

type ChangePasswordReq struct {
	Password        string `json:"password" binding:"required,min=8,max=30"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type UpdateUserReq struct {
	FirstName   *string    `json:"first_name" example:"John"`
	LastName    *string    `json:"last_name" example:"Doe"`
	Phone      *string    `json:"phone" example:"1234567890"`
	CircleStyle *circlestyle.CircleStyleEnum `json:"circle_style"`
}
