package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user/internal/config"
	"user/internal/user"
	"user/internal/user/events"
	"user/pkg/database"
	"user/pkg/event"

	"github.com/caarlos0/env/v11"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	// init database
	err = database.New(context.Background(), cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	// migrate
	err = database.Migrate(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	// rabbitmq
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

	err = ch.ExchangeDeclare(
		event.ExchangeName, // name
		"topic",            // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	var logger *zap.Logger
	var mode string

	switch cfg.Env {
	case "prod":
		mode = gin.ReleaseMode
		l, _ := zap.NewProduction()
		logger = l
	default:
		mode = gin.DebugMode
		l, _ := zap.NewDevelopment()
		logger = l
	}

	gin.SetMode(mode)

	r := gin.New()
	r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
	}))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(cors.Default())

	// Init Event Domain
	// inventoryEvent := inventory.NewInventoryEvent(ch)
	// go inventoryEvent.SubscribeTransactionPaid()
	// go inventoryEvent.SubscribeProductIncreased()

	userEvent := user.NewUserEvent(ch)
	surveyEvent := events.NewSurveyEvent(ch)
	go surveyEvent.SubscribeSurvey()
	go userEvent.SubscribeUser()

	// // Init Router
	userHandler := user.NewHandler(cfg, ch)
	userRouter := user.NewRouter(userHandler, r.RouterGroup)
	userRouter.Register()

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "user service is working")
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	<-ctx.Done()

	log.Println("Server exiting")
}
