package responden

import (
	"api/handler/responden"

	"github.com/gin-gonic/gin"
)

func RespondenRoutes(r *gin.Engine, h *responden.Handler) {
	api := r.Group("/api")
	{
		api.GET("/responden", h.GetRespondenService)
	}
}
