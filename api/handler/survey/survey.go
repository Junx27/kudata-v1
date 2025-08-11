package survey

import (
	"api/config"
	"api/event"
	"api/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
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

func (h *Handler) GetAllSurvey(c *gin.Context) {

	url := fmt.Sprintf("%s%s", h.BaseURL, "/survey")
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
	c.JSON(http.StatusOK, surveys)
}
func (h *Handler) GetAllSurveyCategories(c *gin.Context) {

	url := fmt.Sprintf("%s%s", h.BaseURL, "/categories")
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
	c.JSON(http.StatusOK, surveys)
}

func (h *Handler) CreateSurvey(c *gin.Context) {
	var req model.CreateSurveyRequest
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

	objectName := fmt.Sprintf("survey/%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))

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

	survey := model.MessageSurvey{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Image:       objectName,
		CategoryID:  req.CategoryID,
	}

	surveyData, err := json.Marshal(survey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unmarshal error"})
		return
	}
	err = event.Publisher(h.ch, "create.survey", surveyData)
	c.JSON(http.StatusCreated, gin.H{"message": "Survey created successfully"})
}

func (h *Handler) UpdateSurvey(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey ID"})
		return
	}

	var req model.CreateSurveyRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectName := ""
	file, err := c.FormFile("image")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image"})
			return
		}
		defer src.Close()

		objectName = fmt.Sprintf("survey/%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))

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
	}

	survey := model.MessageSurvey{
		ID:          id,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Image:       objectName,
		CategoryID:  req.CategoryID,
	}

	surveyData, err := json.Marshal(survey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Marshal error"})
		return
	}

	if err := event.Publisher(h.ch, "update.survey", surveyData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish update event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Survey updated successfully"})
}

func (h *Handler) DeleteSurvey(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid survey ID"})
		return
	}

	survey := model.MessageSurvey{
		ID: id,
	}
	surveyData, err := json.Marshal(survey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Marshal error"})
		return
	}

	if err := event.Publisher(h.ch, "delete.survey", surveyData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Survey deleted successfully"})
}
