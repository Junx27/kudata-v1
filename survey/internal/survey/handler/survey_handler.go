package handler

import (
	"context"
	"net/http"
	"strconv"

	"survey/internal/config"
	"survey/internal/survey/repository"

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

func (h *Handler) GetSurveyByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey ID"})
		return
	}

	survey, err := repository.GetSurveyByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get survey"})
		return
	}
	if survey == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
		return
	}

	c.JSON(http.StatusOK, survey)
}

func (h *Handler) GetAllSurvey(c *gin.Context) {
	// Ambil query param
	categoryIDStr := c.Query("category_id")
	name := c.Query("name")

	// Parse category_id jika ada
	var categoryID int
	if categoryIDStr != "" {
		var err error
		categoryID, err = strconv.Atoi(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
			return
		}
	}

	// Panggil service dengan filter
	surveys, err := repository.GetAllSurveys(context.Background(), categoryID, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get surveys"})
		return
	}

	if len(surveys) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No surveys found"})
		return
	}

	c.JSON(http.StatusOK, surveys)
}

func (h *Handler) GetAllCategories(c *gin.Context) {
	categories, err := repository.GetAllCategories(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get categories"})
		return
	}
	if len(categories) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No categories found"})
		return
	}

	c.JSON(http.StatusOK, categories)
}
