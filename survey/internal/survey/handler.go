package survey

import (
	"context"
	"net/http"
	"strconv"
	"survey/internal/config"

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

type createSurveyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (h *Handler) CreateSurvey(c *gin.Context) {
	var req createSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	survey := SurveyInput{
		Name:        req.Name,
		Description: req.Description,
	}

	err := StoreSurvey(context.Background(), survey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store survey"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Survey created successfully"})
}

func (h *Handler) GetSurveyByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey id"})
		return
	}

	survey, err := GetSurveyByID(context.Background(), id)
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

	surveys, err := GetAllSurveys(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get surveys"})
		return
	}
	if surveys == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
		return
	}

	c.JSON(http.StatusOK, surveys)
}
