package router

import (
	"backend/config"
	"backend/internal/constants/enum/jwtpurpose"
	"backend/internal/controller/http/middleware"
	controller "backend/internal/controller/http/v1/controller/user"
	userService "backend/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

type UserServiceSet struct {
	UserAuthService userService.UserAuthService
}

func NewUserRouter(
	router *gin.RouterGroup,
	cfg *config.Config,
	serviceSet *UserServiceSet,
) {
	// New controller
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
	// controller
	{
		privateGroup.POST("/logout", uAuthCtrl.Logout)
	}
}
