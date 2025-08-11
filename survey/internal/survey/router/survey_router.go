package router

import (
	"survey/internal/survey/handler"

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
	r.group.GET("/survey/:id", r.handler.GetSurveyByID)
	r.group.GET("/survey", r.handler.GetAllSurvey)
	r.group.GET("/categories", r.handler.GetAllCategories)

}
