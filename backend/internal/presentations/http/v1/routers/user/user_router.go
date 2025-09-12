package user

import (
	"backend/global"
	"backend/internal/constants/enums/jwtpurpose"
	wireUser "backend/internal/infrastructures/wire/user"
	"backend/internal/presentations/http/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRouter struct {
}

func (u *UserRouter) InitUserRouter(
	router *gin.RouterGroup,
	db *gorm.DB,
) {
	// config
	cfg := global.Config

	uAuthCtrl, _ := wireUser.InitUserRouterHandler(db)
	// ====== Main user group ======
	userGroup := router.Group("/users")

	// ====== Public group ======
	publicGroup := userGroup.Group("")
	{
		publicGroup.POST("/refresh-token", uAuthCtrl.RefreshToken)
		publicGroup.POST("/login", uAuthCtrl.Login)
		publicGroup.POST("/forgot-password", uAuthCtrl.ForgotPassword)
	}

	// Register
	registerGroup := publicGroup.Group("/register")
	{
		registerGroup.POST("/", uAuthCtrl.Register)
		registerGroup.POST("/resend-email", uAuthCtrl.ResendEmailVerifyRegister)
		registerGroup.POST("/verify",
			middlewares.AuthHeader([]byte(cfg.JWT.RegisterTokenKey), jwtpurpose.Register),
			uAuthCtrl.VerifyRegister,
		)
	}

	// ====== Private group (using access token) ======
	privateGroup := userGroup.Group("")
	// middlewares
	privateGroup.Use(middlewares.AuthHeader([]byte(cfg.JWT.AccessTokenKey), jwtpurpose.Access))
	// presentations
	{
		privateGroup.POST("/logout", uAuthCtrl.Logout)
		privateGroup.POST("/change-password", uAuthCtrl.ChangePassword)
	}
}
