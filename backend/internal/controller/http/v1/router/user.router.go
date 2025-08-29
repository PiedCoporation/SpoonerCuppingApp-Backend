package router

import (
	"backend/config"
	"backend/internal/constants/enum/jwtpurpose"
	"backend/internal/controller/http/middleware"
	controller "backend/internal/controller/http/v1/controller/user"
	userService "backend/internal/service/user"

	"github.com/gin-gonic/gin"
)

type UserServiceSet struct {
	authService userService.UserAuthService
}

func NewUserRouter(
	router *gin.RouterGroup,
	cfg *config.Config,
	serviceSet *UserServiceSet,
) {
	// New controller
	uAuthCtrl := controller.NewUserAuthController(serviceSet.authService)

	// ====== Main user group ======
	userGroup := router.Group("/user")

	// ====== Public group ======

	// Register
	publicGroup := userGroup.Group("")
	registerGroup := publicGroup.Group("/register")
	{
		registerGroup.POST("/", uAuthCtrl.Register)
		registerGroup.POST("/resend-email", uAuthCtrl.ResendEmailVerifyRegister)
		registerGroup.POST("/verify",
			middleware.AuthQuery([]byte(cfg.JWT.RegisterTokenKey), jwtpurpose.Register),
			uAuthCtrl.VerifyRegister,
		)
	}

	// Login
	loginGroup := publicGroup.Group("/login")
	{
		loginGroup.POST("/", uAuthCtrl.Login)
		loginGroup.POST("/verify",
			middleware.AuthQuery([]byte(cfg.JWT.LoginTokenKey), jwtpurpose.Login),
			uAuthCtrl.VerifyLogin,
		)
	}

	// ====== Private group (using access token) ======
	privateGroup := userGroup.Group("")
	// middleware
	privateGroup.Use(middleware.AuthHeader([]byte(cfg.JWT.AccessTokenKey), jwtpurpose.Access))
	// controller
	{
		privateGroup.POST("/logout", uAuthCtrl.Logout)
		privateGroup.POST("/refresh-token", uAuthCtrl.RefreshToken)
	}
}
