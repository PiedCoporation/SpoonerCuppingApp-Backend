package routers

import (
	"backend/internal/presentations/http/v1/routers/event"
	"backend/internal/presentations/http/v1/routers/post"
	"backend/internal/presentations/http/v1/routers/sample"
	"backend/internal/presentations/http/v1/routers/user"
)

type RouterGroup struct {
	User  user.RouterGroup
	Event event.RouterGroup
	Post  post.RouterGroup
	Sample sample.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
