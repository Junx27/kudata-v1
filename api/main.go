package main

import (
	"fmt"
	"log"
	"net/http"

	"api/config"
	"api/event"
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

	config.AppConfig = cfg
	config.InitMinio()

	// connect to RabbitMQ
	conn, err := event.New(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ch.Close()

	r := gin.Default()

	userHandler := user.NewHandler(cfg, ch, cfg.BaseURLUser)
	userRoute.UserRoutes(r, &userHandler)

	surveyHandler := survey.NewHandler(cfg, ch, cfg.BaseURLSurvey)
	surveyRoute.SurveyRoutes(r, &surveyHandler)

	respondenHandler := &responden.Handler{
		BaseURL: cfg.BaseURLResponden,
	}
	respondenRoute.RespondenRoutes(r, respondenHandler)

	paymentHandler := &payment.Handler{
		BaseURL: cfg.BaseURLPayment,
	}
	paymentRoute.PaymentRoutes(r, paymentHandler)

	r.GET("/api/status-service", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "API service is working"})
	})

	if err := r.Run(":8004"); err != nil {
		log.Fatalf("Error starting API Gateway: %v", err)
	}
}
