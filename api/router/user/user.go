package user

import (
	"api/handler/user"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, h *user.Handler) {
	api := r.Group("/api")
	{
		api.GET("/users", h.GetAllUsers)
	}
}
