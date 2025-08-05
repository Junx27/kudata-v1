package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	BaseURL string
}

func (h *Handler) decodeResponseBody(resp *http.Response) ([]map[string]interface{}, error) {
	var surveys []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&surveys); err != nil {
		return nil, fmt.Errorf("Error decoding response body: %v", err)
	}
	return surveys, nil
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

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
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

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    user,
	})
}
