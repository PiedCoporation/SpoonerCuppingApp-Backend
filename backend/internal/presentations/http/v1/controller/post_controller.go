package controller

import (
	"backend/internal/constants/errorcode"
	"backend/internal/usecases/abstractions"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	postService abstractions.IPostService
}

func NewPostController(
	postService abstractions.IPostService,
) *PostController {
	return &PostController{
		postService: postService,
	}
}

// GetPosts godoc
// @Summary Get paginated posts
// @Description Retrieve a paginated list of posts with optional search functionality
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page_size query int false "Number of posts per page (default: 10, max: 100)" default(10) minimum(1) maximum(100)
// @Param page_number query int false "Page number (default: 1)" default(1) minimum(1)
// @Param search_term query string false "Search term to filter posts by title"
// @Success 200 {object} post.PostPageResult "Paginated list of posts"
// @Failure 400 {object} controller.ErrorResponse "Bad request - invalid parameters"
// @Failure 401 {object} controller.ErrorResponse "Unauthorized - invalid or missing token"
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
