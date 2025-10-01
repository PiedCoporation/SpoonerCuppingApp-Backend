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

type UserController struct {
	userSettingService abstractions.IUserSettingService
}

func NewUserController(
	userSettingService abstractions.IUserSettingService,
) *UserController {
	return &UserController{
		userSettingService: userSettingService,
	}
}

// GetUserSetting godoc
// @Summary Get user setting
// @Description Retrieve a user setting by its ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} user.UserSettingRes "User setting"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized"
// @Failure 403 {object} controller.ErrorResponse "Forbidden"
// @Failure 404 {object} controller.ErrorResponse "User not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /users/settings [get]
func (uc *UserController) GetUserSetting(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}

	ctx := c.Request.Context()
	setting, err := uc.userSettingService.GetByID(ctx, userID.(uuid.UUID))
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, setting)
}

// UpdateUserSetting godoc
// @Summary Update an existing user setting
// @Description Update setting's circleStyle
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body user.UpdateUserSettingReq true "Update payload"
// @Success 200 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse "Bad request"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized"
// @Failure 403 {object} controller.ErrorResponse "Forbidden"
// @Failure 404 {object} controller.ErrorResponse "Not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /users/settings [patch]
func (uc *UserController) UpdateUserSetting(c *gin.Context) {
	var req user.UpdateUserSettingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}

	ctx := c.Request.Context()
	if err := uc.userSettingService.Update(ctx, userID.(uuid.UUID), req); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "setting updated"})
}
