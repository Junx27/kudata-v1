package router

import (
	"user/internal/user/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	handler handler.Handler
	group   gin.RouterGroup
}

func NewRouter(handler handler.Handler, group gin.RouterGroup) Router {
	return Router{
		handler: handler,
		group:   group,
	}
}

func (r *Router) Register() {
	r.group.GET("/user/:id", r.handler.GetUserByID)
	r.group.POST("/user", r.handler.CreateUser)
	r.group.GET("/user", r.handler.GetAllUsers)
}
