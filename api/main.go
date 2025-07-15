package main

import (
	"fmt"
	"log"

	"api/config"
	"api/handler/payment"
	"api/handler/responden"
	"api/handler/survey"
	"api/handler/user"
	paymentRoute "api/router/payment"
	respondenRoute "api/router/responden"
	surveyRoute "api/router/survey"
	userRoute "api/router/user"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	r := gin.Default()

	userHandler := &user.Handler{
		BaseURL: cfg.BaseURLUser,
	}
	userRoute.UserRoutes(r, userHandler)

	surveyHandler := &survey.Handler{
		BaseURL: cfg.BaseURLSurvey,
	}
	surveyRoute.SurveyRoutes(r, surveyHandler)

	respondenHandler := &responden.Handler{
		BaseURL: cfg.BaseURLResponden,
	}
	respondenRoute.RespondenRoutes(r, respondenHandler)

	paymentHandler := &payment.Handler{
		BaseURL: cfg.BaseURLPayment,
	}
	paymentRoute.PaymentRoutes(r, paymentHandler)

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "API Gateway is working")
	})

	if err := r.Run(":8004"); err != nil {
		log.Fatalf("Error starting API Gateway: %v", err)
	}
}
