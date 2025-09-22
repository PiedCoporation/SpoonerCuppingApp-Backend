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

type PostController struct {
	postService    abstractions.IPostService
	commentService abstractions.ICommentService
}

func NewPostController(
	postService abstractions.IPostService,
	commentService abstractions.ICommentService,
) *PostController {
	return &PostController{
		postService:    postService,
		commentService: commentService,
	}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post with optional images
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body post.CreatePostReq true "Post payload"
// @Success 201 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse "Bad request"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /posts [post]
func (pc *PostController) CreatePost(c *gin.Context) {
	var req post.CreatePostReq
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
	if err := pc.postService.Create(ctx, userID.(uuid.UUID), req); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "post created"})
}

// GetPosts godoc
// @Summary Get paginated posts
// @Description Retrieve a paginated list of posts with optional search functionality
// @Tags posts
// @Accept json
// @Produce json
// @Param page_size query int false "Number of posts per page (default: 10, max: 100)" default(10) minimum(1) maximum(100)
// @Param page_number query int false "Page number (default: 1)" default(1) minimum(1)
// @Param search_term query string false "Search term to filter posts by title"
// @Success 200 {object} post.PostPageResult "Paginated list of posts"
// @Failure 400 {object} controller.ErrorResponse "Bad request - invalid parameters"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /posts [get]
func (pc *PostController) GetPosts(c *gin.Context) {
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil {
		pageSize = 10
	}

	pageNumber, err := strconv.Atoi(c.Query("page_number"))
	if err != nil {
		pageNumber = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if pageNumber < 1 {
		pageNumber = 1
	}

	searchTerm := c.Query("search_term")
	ctx := c.Request.Context()
	posts, err := pc.postService.GetAll(ctx, pageSize, pageNumber, searchTerm)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

// GetPostByID godoc
// @Summary Get post by ID
// @Description Retrieve a specific post by its ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID" format(uuid)
// @Success 200 {object} post.GetPostByIdResponse "Post details"
// @Failure 400 {object} controller.ErrorResponse "Bad request - invalid parameters"
// @Failure 404 {object} controller.ErrorResponse "Post not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /posts/{id} [get]
func (pc *PostController) GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	ctx := c.Request.Context()
	posts, err := pc.postService.GetByID(ctx, id)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (pc *PostController) GetDirectChildrenComments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	ctx := c.Request.Context()
	comments, err := pc.commentService.GetDirectChildren(ctx, id, true)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, comments)
}

// UpdatePost godoc
// @Summary Update an existing post
// @Description Update a post's title, content, event, and images
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID" format(uuid)
// @Param body body post.UpdatePostReq true "Post payload"
// @Success 200 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse "Bad request"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized"
// @Failure 403 {object} controller.ErrorResponse "Forbidden"
// @Failure 404 {object} controller.ErrorResponse "Not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /posts/{id} [patch]
func (pc *PostController) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	var req post.UpdatePostReq
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
	if err := pc.postService.Update(ctx, userID.(uuid.UUID), id, req); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post updated"})
}

// DeletePost godoc
// @Summary Delete an existing post
// @Description Soft delete post by id
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID" format(uuid)
// @Success 200 {object} controller.MessageResponse
// @Failure 400 {object} controller.ErrorResponse "Bad request"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized"
// @Failure 403 {object} controller.ErrorResponse "Forbidden"
// @Failure 404 {object} controller.ErrorResponse "Not found"
// @Failure 500 {object} controller.ErrorResponse "Internal server error"
// @Router /posts/{id} [delete]
func (pc *PostController) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}

	ctx := c.Request.Context()
	if err := pc.postService.Delete(ctx, userID.(uuid.UUID), id); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}
