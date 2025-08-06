package survey

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"survey/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
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
	Name        string `form:"name" binding:"required"`
	Price       int    `form:"price" binding:"required"`
	Description string `form:"description" binding:"required"`
}

func (h *Handler) CreateSurvey(c *gin.Context) {
	var req createSurveyRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil file image
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image"})
		return
	}
	defer src.Close()

	// Generate unique object name
	objectName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))

	// Upload ke MinIO
	_, err = config.MinioClient.PutObject(context.Background(),
		h.cfg.MinioBucket,
		objectName,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	imageURL := objectName

	survey := SurveyInput{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	}

	err = StoreSurvey(context.Background(), survey, imageURL)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey ID"})
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
	if len(surveys) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No surveys found"})
		return
	}

	c.JSON(http.StatusOK, surveys)
}
