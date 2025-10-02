package controller

import (
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/post"
	"backend/internal/usecases/abstractions"
	"backend/pkg/utils/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostCommentController struct {
	postCmtService abstractions.IPostCommentService
}

func NewPostCommentController(
	postCmtService abstractions.IPostCommentService,
) *PostCommentController {
	return &PostCommentController{
		postCmtService: postCmtService,
	}
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a new comment with post id and optional parent id
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID" format(uuid)
// @Param body body post.CreatePostCommentReq true "PostComment payload"
// @Success 201 {object} controller.IdMessageResponse
// @Failure 400 {object} controller.ErrorResponse "Bad request"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /posts/{id}/comments [post]
func (cc *PostCommentController) CreateComment(c *gin.Context) {
	postIdStr := c.Param("id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		errorcode.JSONError(c, errorcode.ErrInvalidParams)
		return
	}

	var req post.CreatePostCommentReq
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
	id, err := cc.postCmtService.Create(ctx, userID.(uuid.UUID), postId, req)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "comment created",
		"id":      id,
	})
}

// GetRootComments godoc
// @Summary Get root comments in post
// @Description Retrieve a list of root comments by post id
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Post ID" format(uuid)
// @Param order_by_desc query boolean false "Order by desc"
// @Success 200 {array} post.PostCommentViewRes "Root comments"
// @Failure 400 {object} controller.ErrorResponse "Bad request - invalid parameters"
// @Failure 404 {object} controller.ErrorResponse "Post not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /posts/{id}/comments [get]
func (cc *PostCommentController) GetRootComments(c *gin.Context) {
	postIdStr := c.Param("id")
	postId, err := uuid.Parse(postIdStr)
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
	comments, err := cc.postCmtService.GetRootCommentsByPostID(ctx, postId, orderByDesc)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, comments)
}

// GetDirectChildren godoc
// @Summary Get direct children in comment
// @Description Retrieve a list of direct children by comment id
// @Tags comments
// @Accept json
// @Produce json
// @Param commentId path string true "PostComment ID" format(uuid)
// @Param order_by_desc query boolean false "Order by desc"
// @Success 200 {array} post.PostCommentViewRes "Direct children"
// @Failure 400 {object} controller.ErrorResponse "Bad request - invalid parameters"
// @Failure 404 {object} controller.ErrorResponse "PostComment not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /posts/{id}/comments/{commentId}/children [get]
func (cc *PostCommentController) GetDirectChildren(c *gin.Context) {
	idStr := c.Param("commentId")
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
	comments, err := cc.postCmtService.GetDirectChildren(ctx, id, orderByDesc)
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
// @Param commentId path string true "PostComment ID" format(uuid)
// @Param body body post.UpdatePostCommentReq true "PostComment payload"
// @Success 200 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse "Bad request"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized"
// @Failure 403 {object} controller.ErrorResponse "Forbidden"
// @Failure 404 {object} controller.ErrorResponse "Not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /posts/{id}/comments/{commentId}/ [put]
func (cc *PostCommentController) UpdateComment(c *gin.Context) {
	idStr := c.Param("commentId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errorcode.JSONError(c, errorcode.ErrInvalidParams)
		return
	}

	var req post.UpdatePostCommentReq
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
	if err := cc.postCmtService.Update(ctx, userID.(uuid.UUID), id, req); err != nil {
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
// @Param commentId path string true "PostComment ID" format(uuid)
// @Success 200 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse "Bad request"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized"
// @Failure 403 {object} controller.ErrorResponse "Forbidden"
// @Failure 404 {object} controller.ErrorResponse "Not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /posts/{id}/comments/{commentId}/ [delete]
func (cc *PostCommentController) DeleteComment(c *gin.Context) {
	idStr := c.Param("commentId")
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
	if err := cc.postCmtService.Delete(ctx, userID.(uuid.UUID), id); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}
