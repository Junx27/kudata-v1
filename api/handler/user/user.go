package user

import (
	"api/config"
	"api/event"
	"api/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Handler struct {
	cfg     config.Config
	ch      *amqp.Channel
	BaseURL string
}

func NewHandler(cfg config.Config, ch *amqp.Channel, baseURL string) Handler {
	return Handler{
		cfg:     cfg,
		ch:      ch,
		BaseURL: baseURL,
	}
}

func (h *Handler) decodeResponseBody(resp *http.Response) ([]map[string]interface{}, error) {
	var surveys []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&surveys); err != nil {
		return nil, fmt.Errorf("Error decoding response body: %v", err)
	}
	return surveys, nil
}
func (h *Handler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.MessageUser{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	userData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unmarshal error"})
		return
	}
	err = event.Publisher(h.ch, "create.user", userData)
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})

}
func (h *Handler) GetAllUsers(c *gin.Context) {

	url := fmt.Sprintf("%s%s", h.BaseURL, "/user")
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch surveys"})
		return
	}
	users, err := h.decodeResponseBody(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    users,
	})
}

func (h *Handler) GetUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	url := fmt.Sprintf("%s/user/%d", h.BaseURL, id)
	log.Println("user url:", url)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "user service returned error"})
		return
	}

	var user model.UseResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    user,
	})
}
