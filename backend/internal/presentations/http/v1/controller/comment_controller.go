package controller

import (
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/comment"
	"backend/internal/usecases/abstractions"
	"backend/pkg/utils/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentController struct {
	commentService abstractions.ICommentService
}

func NewCommentController(
	commentService abstractions.ICommentService,
) *CommentController {
	return &CommentController{
		commentService: commentService,
	}
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a new comment with post id and optional parent id
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body comment.CreateCommentReq true "Comment payload"
// @Success 201 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse "Bad request"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /comments [post]
func (cc *CommentController) CreateComment(c *gin.Context) {
	var req comment.CreateCommentReq
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
	if err := cc.commentService.Create(ctx, userID.(uuid.UUID), req); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "comment created"})
}

// GetDirectChildren godoc
// @Summary Get direct children in comment
// @Description Retrieve a list of direct children by comment id
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID" format(uuid)
// @Param order_by_desc query boolean false "Order by desc"
// @Success 200 {array} comment.CommentViewRes "Direct children"
// @Failure 400 {object} controller.ErrorResponse "Bad request - invalid parameters"
// @Failure 404 {object} controller.ErrorResponse "Comment not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /comments/{id}/children [get]
func (cc *CommentController) GetDirectChildren(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errorcode.JSONError(c, errorcode.ErrInvalidParams)
		return
	}

	orderByDescStr := c.DefaultQuery("order_by_desc", "true")
	orderByDesc, err := strconv.ParseBool(orderByDescStr)
	if err != nil {
		orderByDesc = true
	}

	ctx := c.Request.Context()
	comments, err := cc.commentService.GetDirectChildren(ctx, id, orderByDesc)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, comments)
}

// UpdateComment godoc
// @Summary Update an existing comment
// @Description Update a comment's content
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Comment ID" format(uuid)
// @Param body body comment.UpdateCommentReq true "Comment payload"
// @Success 200 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse "Bad request"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized"
// @Failure 403 {object} controller.ErrorResponse "Forbidden"
// @Failure 404 {object} controller.ErrorResponse "Not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /comments/{id} [put]
func (cc *CommentController) UpdateComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errorcode.JSONError(c, errorcode.ErrInvalidParams)
		return
	}

	var req comment.UpdateCommentReq
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
	if err := cc.commentService.Update(ctx, userID.(uuid.UUID), id, req); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment updated"})
}

// DeletePost godoc
// @Summary Delete an existing comment
// @Description Soft delete comment by id
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Comment ID" format(uuid)
// @Success 200 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse "Bad request"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized"
// @Failure 403 {object} controller.ErrorResponse "Forbidden"
// @Failure 404 {object} controller.ErrorResponse "Not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /comments/{id} [delete]
func (cc *CommentController) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errorcode.JSONError(c, errorcode.ErrInvalidParams)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}

	ctx := c.Request.Context()
	if err := cc.commentService.Delete(ctx, userID.(uuid.UUID), id); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}
