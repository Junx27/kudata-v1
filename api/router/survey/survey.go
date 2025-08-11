package survey

import (
	"api/handler/survey"

	"github.com/gin-gonic/gin"
)

func SurveyRoutes(r *gin.Engine, h *survey.Handler) {
	api := r.Group("/api")
	{
		api.GET("/surveys", h.GetAllSurvey)
		api.POST("/surveys", h.CreateSurvey)
		api.PUT("/surveys/:id", h.UpdateSurvey)
		api.DELETE("/surveys/:id", h.DeleteSurvey)
		api.GET("/categories", h.GetAllSurveyCategories)
	}
}
