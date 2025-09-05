package router

import (
	"backend/global"
	"backend/internal/constants/enum/jwtpurpose"
	"backend/internal/presentations/http/middleware"
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
	}

	// Register
	registerGroup := publicGroup.Group("/register")
	{
		registerGroup.POST("/", uAuthCtrl.Register)
		registerGroup.POST("/resend-email", uAuthCtrl.ResendEmailVerifyRegister)
		registerGroup.GET("/verify",
			middleware.AuthQuery([]byte(cfg.JWT.RegisterTokenKey), jwtpurpose.Register),
			uAuthCtrl.VerifyRegister,
		)
	}

	// Login
	loginGroup := publicGroup.Group("/login")
	{
		loginGroup.POST("/", uAuthCtrl.Login)
		loginGroup.GET("/verify",
			middleware.AuthQuery([]byte(cfg.JWT.LoginTokenKey), jwtpurpose.Login),
			uAuthCtrl.VerifyLogin,
		)
	}

	// ====== Private group (using access token) ======
	privateGroup := userGroup.Group("")
	// middleware
	privateGroup.Use(middleware.AuthHeader([]byte(cfg.JWT.AccessTokenKey), jwtpurpose.Access))
	// presentations
	{
		privateGroup.POST("/logout", uAuthCtrl.Logout)
	}
}
