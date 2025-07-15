package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"survey/internal/config"
	"survey/internal/survey"
	"survey/internal/survey/events"
	"survey/pkg/database"
	"survey/pkg/event"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
)

var (
	// WebSocket upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// Connected WebSocket clients
	wsClients   = make(map[*websocket.Conn]bool)
	wsClientsMu sync.Mutex
)

// Broadcast message to all connected clients
func broadcastWebSocketMessage(msg string) {
	wsClientsMu.Lock()
	defer wsClientsMu.Unlock()

	for conn := range wsClients {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Println("WebSocket write error:", err)
			conn.Close()
			delete(wsClients, conn)
		}
	}
}

func main() {
	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	// init database
	if err := database.New(context.Background(), cfg); err != nil {
		log.Fatal(err.Error())
	}

	// migrate
	if err := database.Migrate(cfg); err != nil {
		log.Fatal(err.Error())
	}

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

	// Declare exchange
	err = ch.ExchangeDeclare(
		event.ExchangeName, // name
		"topic",            // type
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	// setup logger
	var logger *zap.Logger
	var mode string
	switch cfg.Env {
	case "prod":
		mode = gin.ReleaseMode
		logger, _ = zap.NewProduction()
	default:
		mode = gin.DebugMode
		logger, _ = zap.NewDevelopment()
	}
	gin.SetMode(mode)

	r := gin.New()
	r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
	}))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(cors.Default())

	// WebSocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}

		wsClientsMu.Lock()
		wsClients[conn] = true
		wsClientsMu.Unlock()

		log.Println("WebSocket client connected")

		// Listen for disconnect
		go func() {
			defer func() {
				wsClientsMu.Lock()
				delete(wsClients, conn)
				wsClientsMu.Unlock()
				conn.Close()
				log.Println("WebSocket client disconnected")
			}()
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					break
				}
			}
		}()
	})

	// survey events
	surveyEvent := survey.NewSurveyEvent(ch)
	userEvent := events.NewUserEvent(ch)
	go userEvent.SubscribeUser()
	go surveyEvent.SubscribeSurvey(broadcastWebSocketMessage) // pass function

	// survey route
	surveyHandler := survey.NewHandler(cfg, ch)
	surveyRouter := survey.NewRouter(surveyHandler, r.RouterGroup)
	surveyRouter.Register()

	// plain HTTP GET to show latest message
	r.GET("/", func(ctx *gin.Context) {
		survey.MessageMu.RLock()
		defer survey.MessageMu.RUnlock()

		msg := survey.LatestMessage
		if msg == "" {
			ctx.String(http.StatusOK, "No message received yet")
		} else {
			ctx.String(http.StatusOK, "Last message: %s", msg)
		}
	})

	r.POST("/payment", func(c *gin.Context) {
		var payload map[string]interface{}

		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize JSON"})
			return
		}

		broadcastWebSocketMessage(string(jsonBytes))

		c.JSON(http.StatusOK, gin.H{
			"status":  "broadcasted",
			"payload": payload,
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r.Handler(),
	}

	// run server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
