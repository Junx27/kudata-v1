package survey

import "github.com/gin-gonic/gin"

type Router struct {
	handler Handler
	group   gin.RouterGroup
}

func NewRouter(handler Handler, group gin.RouterGroup) Router {
	return Router{
		handler: handler,
		group:   group,
	}
}

func (r *Router) Register() {
	r.group.GET("/survey/:id", r.handler.GetSurveyByID)
	r.group.POST("/survey", r.handler.CreateSurvey)
	r.group.GET("/survey", r.handler.GetAllSurvey)
}
