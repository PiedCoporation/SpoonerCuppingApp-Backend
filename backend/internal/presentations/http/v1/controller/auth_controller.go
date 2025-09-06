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
	auth abstractions.UserAuthService
}

func NewUserAuthController(
	auth abstractions.UserAuthService,
) *UserAuthController {
	return &UserAuthController{
		auth: auth,
	}
}

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
	}

	if err := uc.auth.Register(ctx, vo); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Please check your email to verify account.",
	})
}

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

func (uc *UserAuthController) VerifyRegister(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID in token"})
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

func (uc *UserAuthController) Login(c *gin.Context) {
	var req user.LoginUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	ctx := c.Request.Context()
	vo := user.LoginUserVO{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := uc.auth.Login(ctx, vo); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Please check your email to login account.",
	})
}

func (uc *UserAuthController) VerifyLogin(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID in token"})
		return
	}

	ctx := c.Request.Context()
	accessToken, refreshToken, err := uc.auth.VerifyLogin(ctx, userID.(uuid.UUID))
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"data": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (uc *UserAuthController) Logout(c *gin.Context) {
	var req user.LogoutUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	// get userID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID in token"})
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
