package payment

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func (h *Handler) GetPaymentService(c *gin.Context) {

	url := fmt.Sprintf("%s", h.BaseURL)
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
	surveys, err := h.decodeResponseBody(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    surveys,
	})
}
