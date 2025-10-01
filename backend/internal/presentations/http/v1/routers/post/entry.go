package post

import "backend/internal/presentations/http/v1/routers/post/comment"

type RouterGroup struct {
	PostRouter
	Comment comment.RouterGroup
}

var PostRouterGroup = new(RouterGroup)
