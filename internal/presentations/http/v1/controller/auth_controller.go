package controller

import (
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/user"
	"backend/internal/usecases/abstractions"
	"backend/pkg/utils/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserAuthController struct {
	auth abstractions.IUserAuthService
}

func NewUserAuthController(
	auth abstractions.IUserAuthService,
) *UserAuthController {
	return &UserAuthController{
		auth: auth,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a user account
// @Tags auth
// @Accept json
// @Produce json
// @Param body body user.RegisterUserReq true "Register payload"
// @Success 200 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse
// @Failure 500 {object} controller.ErrorResponse
// @Router /users/register [post]
func (uc *UserAuthController) Register(c *gin.Context) {
	var req user.RegisterUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	ctx := c.Request.Context()
	vo := user.RegisterUserVO{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
	}

	if err := uc.auth.Register(ctx, vo); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Please check your email to verify account.",
	})
}

// ResendEmailVerifyRegister godoc
// @Summary Resend verification email for register
// @Description Resend email verification code for registration
// @Tags auth
// @Accept json
// @Produce json
// @Param body body user.ResendEmailReq true "Email payload"
// @Success 200 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse
// @Failure 500 {object} controller.ErrorResponse
// @Router /users/register/resend-email [post]
func (uc *UserAuthController) ResendEmailVerifyRegister(c *gin.Context) {
	var req user.ResendEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	ctx := c.Request.Context()
	email := req.Email

	if err := uc.auth.ResendEmailVerifyRegister(ctx, email); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Please check your email to verify account.",
	})
}

// VerifyRegister godoc
// @Summary Verify registration and issue tokens
// @Description Verify user registration using token and return access/refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} controller.TokenResponse
// @Failure 401 {object} controller.ErrorResponse
// @Failure 500 {object} controller.ErrorResponse
// @Router /users/register/verify [post]
func (uc *UserAuthController) VerifyRegister(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}

	ctx := c.Request.Context()
	accessToken, refreshToken, err := uc.auth.VerifyRegister(ctx, userID.(uuid.UUID))
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "register success",
		"data": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// Login godoc
// @Summary Login
// @Description Login and obtain access/refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param body body user.LoginUserReq true "Login payload"
// @Success 200 {object} controller.TokenResponse
// @Failure 400 {object} controller.ErrorResponse
// @Failure 401 {object} controller.ErrorResponse
// @Failure 500 {object} controller.ErrorResponse
// @Router /users/login [post]
func (uc *UserAuthController) Login(c *gin.Context) {
	var req user.LoginUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	ctx := c.Request.Context()
	// vo := user.LoginUserReq{
	// 	Email:    req.Email,
	// 	Password: req.Password,
	// }

	result := uc.auth.Login(ctx, req)
	// if result.IsFailure {
	// 	errorcode.JSONError(c, err)
	// 	return
	// }

	c.JSON(http.StatusOK, result)
}

// Logout godoc
// @Summary Logout
// @Description Revoke refresh token and logout
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body user.LogoutUserReq true "Logout payload"
// @Success 200 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse
// @Failure 401 {object} controller.ErrorResponse
// @Failure 500 {object} controller.ErrorResponse
// @Router /users/logout [post]
func (uc *UserAuthController) Logout(c *gin.Context) {
	var req user.LogoutUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	// get userID from middlewares
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}
	ctx := c.Request.Context()

	err := uc.auth.Logout(ctx, userID.(uuid.UUID), req.RefreshToken)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}

// ForgotPassword godoc
// @Summary Send forgot password email
// @Description Send email to reset password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body user.ForgotPasswordReq true "Email payload"
// @Success 200 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse
// @Failure 500 {object} controller.ErrorResponse
// @Router /users/forgot-password [post]
func (uc *UserAuthController) ForgotPassword(c *gin.Context) {
	var req user.ForgotPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	ctx := c.Request.Context()
	email := req.Email

	if err := uc.auth.ForgotPassword(ctx, email); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Please check your email to change password.",
	})
}

// ChangePassword godoc
// @Summary Change password
// @Description Change password for the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body user.ChangePasswordReq true "Change password payload"
// @Success 200 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse
// @Failure 401 {object} controller.ErrorResponse
// @Failure 500 {object} controller.ErrorResponse
// @Router /users/change-password [post]
func (uc *UserAuthController) ChangePassword(c *gin.Context) {
	var req user.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	// get userID from middlewares
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}
	ctx := c.Request.Context()
	vo := user.ChangePasswordVO{
		UserID:   userID.(uuid.UUID),
		Password: req.Password,
	}

	if err := uc.auth.ChangePassword(ctx, vo); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Change password success.",
	})
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Exchange refresh token for new tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param body body user.RefreshTokenReq true "Refresh token payload"
// @Success 200 {object} controller.TokenResponse
// @Failure 400 {object} controller.ErrorResponse
// @Failure 401 {object} controller.ErrorResponse
// @Failure 500 {object} controller.ErrorResponse
// @Router /users/refresh-token [post]
func (uc *UserAuthController) RefreshToken(c *gin.Context) {
	var req user.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	ctx := c.Request.Context()

	accessToken, refreshToken, err := uc.auth.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "refresh token success",
		"data": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (uc *UserAuthController) UserUpdate(c *gin.Context) {}