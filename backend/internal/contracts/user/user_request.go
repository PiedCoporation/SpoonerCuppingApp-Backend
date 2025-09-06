package user

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
