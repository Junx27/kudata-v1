package payment

import (
	"api/handler/payment"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(r *gin.Engine, h *payment.Handler) {
	api := r.Group("/api")
	{
		api.GET("/payment", h.GetPaymentService)
	}
}
