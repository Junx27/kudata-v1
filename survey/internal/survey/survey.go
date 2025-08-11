package survey

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"survey/internal/config"
	"survey/pkg/database"
	"survey/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type SurveyInput struct {
	Name        string `form:"name" binding:"required"`
	Price       int    `form:"price" binding:"required"`
	Description string `form:"description" binding:"required"`
	CategoryID  int    `form:"category_id" binding:"required"`
}

type SurveyResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type CategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func StoreSurvey(ctx context.Context, req SurveyInput, imageURL string) error {
	query := `INSERT INTO surveys (name, image, price, description, category_id) 
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := database.DB.Exec(ctx, query, req.Name, imageURL, req.Price, req.Description, req.CategoryID)
	if err != nil {
		log.Println("Error inserting survey:", err)
		return err
	}
	return nil
}

func CreateSurveyHandler(c *gin.Context) {
	var input SurveyInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}

	imageURL, err := storage.UploadImageToMinio(file, fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
		return
	}

	if err := StoreSurvey(c.Request.Context(), input, imageURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store survey"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "survey created successfully",
		"image":   imageURL,
	})
}

func GetSurveyByID(ctx context.Context, id int) (*SurveyResponse, error) {
	query := `
		SELECT s.id, s.name, s.image, s.price, s.description, c.name AS category
		FROM surveys s
		JOIN categories c ON s.category_id = c.id
		WHERE s.id = $1
	`

	row := database.DB.QueryRow(ctx, query, id)

	var survey SurveyResponse
	err := row.Scan(
		&survey.ID,
		&survey.Name,
		&survey.Image,
		&survey.Price,
		&survey.Description,
		&survey.Category,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Println("Error fetching survey:", err)
		return nil, err
	}

	// Format image URL
	url := config.AppConfig.MinioHost
	if config.AppConfig.MinioHost == "minio:9000" {
		url = "localhost:9000"
	}
	survey.Image = fmt.Sprintf("http://%s/%s/%s",
		url,
		config.AppConfig.MinioBucket,
		survey.Image,
	)

	return &survey, nil
}

func GetAllSurveys(ctx context.Context, categoryID int, name string) ([]SurveyResponse, error) {
	query := `
		SELECT s.id, s.name, s.image, s.price, s.description, c.name AS category
		FROM surveys s
		JOIN categories c ON s.category_id = c.id
		WHERE 1=1
	`

	args := []interface{}{}
	i := 1

	if categoryID != 0 {
		query += fmt.Sprintf(" AND s.category_id = $%d", i)
		args = append(args, categoryID)
		i++
	}

	if name != "" {
		query += fmt.Sprintf(" AND s.name ILIKE $%d", i)
		args = append(args, "%"+name+"%")
		i++
	}

	rows, err := database.DB.Query(ctx, query, args...)
	if err != nil {
		log.Println("Error querying survey:", err)
		return nil, err
	}
	defer rows.Close()

	var surveys []SurveyResponse

	for rows.Next() {
		var survey SurveyResponse
		err := rows.Scan(
			&survey.ID,
			&survey.Name,
			&survey.Image,
			&survey.Price,
			&survey.Description,
			&survey.Category,
		)
		if err != nil {
			log.Println("Error scanning survey:", err)
			return nil, err
		}

		url := config.AppConfig.MinioHost
		if config.AppConfig.MinioHost == "minio:9000" {
			url = "localhost:9000"
		}
		survey.Image = fmt.Sprintf("http://%s/%s/%s",
			url,
			config.AppConfig.MinioBucket,
			survey.Image,
		)

		surveys = append(surveys, survey)
	}

	if err = rows.Err(); err != nil {
		log.Println("Rows iteration error:", err)
		return nil, err
	}

	return surveys, nil
}

func GetAllCategories(ctx context.Context) ([]CategoryResponse, error) {
	query := `SELECT id, name FROM categories`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		log.Println("Error querying category:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryResponse

	for rows.Next() {
		var category CategoryResponse
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			log.Println("Error scanning category:", err)
			return nil, err
		}
		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		log.Println("Rows iteration error:", err)
		return nil, err
	}

	return categories, nil
}
