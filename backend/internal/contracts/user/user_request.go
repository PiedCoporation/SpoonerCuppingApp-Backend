package user

type RegisterUserReq struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
}

type ResendEmailReq struct {
	Email string `json:"email" binding:"required,email"`
}

type LoginUserReq struct {
	Email string `json:"email" binding:"required,email"`
}

type LogoutUserReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
