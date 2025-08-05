package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"user/internal/config"
	"user/pkg/event"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Handler struct {
	cfg config.Config
	ch  *amqp.Channel
}

func NewHandler(cfg config.Config, ch *amqp.Channel) Handler {
	return Handler{
		cfg: cfg,
		ch:  ch,
	}
}

type createUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := UserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := StoreUser(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user"})
		return
	}

	userData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unmarshal error"})
		return
	}

	err = event.Publisher(h.ch, "create.user.success", userData)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	user, err := GetUserByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetAllUsers(c *gin.Context) {

	user, err := GetAllUsers(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
