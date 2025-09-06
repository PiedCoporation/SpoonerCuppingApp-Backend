package router

import (
	"backend/global"
	"backend/internal/constants/enums/jwtpurpose"
	"backend/internal/presentations/http/middlewares"
	"backend/internal/presentations/http/v1/controller"
	userService "backend/internal/usecases/abstractions"

	"github.com/gin-gonic/gin"
)

type UserServiceSet struct {
	UserAuthService userService.UserAuthService
}

func NewUserRouter(
	router *gin.RouterGroup,
	serviceSet *UserServiceSet,
) {
	// config
	cfg := global.Config

	// New presentations
	uAuthCtrl := controller.NewUserAuthController(serviceSet.UserAuthService)

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
