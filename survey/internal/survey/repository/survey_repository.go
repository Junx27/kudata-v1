package repository

import (
	"context"
	"fmt"
	"log"
	"survey/internal/config"
	"survey/internal/survey/model"
	"survey/pkg/database"

	"github.com/jackc/pgx/v5"
)

func StoreSurvey(ctx context.Context, req model.SurveyInput, imageURL string) error {
	query := `INSERT INTO surveys (name, image, price, description, category_id) 
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := database.DB.Exec(ctx, query, req.Name, imageURL, req.Price, req.Description, req.CategoryID)
	if err != nil {
		log.Println("Error inserting survey:", err)
		return err
	}
	return nil
}

func UpdateSurvey(ctx context.Context, id int, req model.SurveyInput, imageURL string) error {
	query := `UPDATE surveys 
			  SET name = $1, image = $2, price = $3, description = $4, category_id = $5
			  WHERE id = $6`
	_, err := database.DB.Exec(ctx, query, req.Name, imageURL, req.Price, req.Description, req.CategoryID, id)
	if err != nil {
		log.Println("Error updating survey:", err)
		return err
	}
	return nil
}

func DeleteSurvey(ctx context.Context, id int) error {
	query := `DELETE FROM surveys WHERE id = $1`
	_, err := database.DB.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error deleting survey:", err)
		return err
	}
	return nil
}

func GetSurveyByID(ctx context.Context, id int) (*model.SurveyResponse, error) {
	query := `
		SELECT s.id, s.name, s.image, s.price, s.description, c.name AS category
		FROM surveys s
		JOIN categories c ON s.category_id = c.id
		WHERE s.id = $1
	`

	row := database.DB.QueryRow(ctx, query, id)

	var survey model.SurveyResponse
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
	survey.Image = fmt.Sprintf("http://%s/%s/%s/%s",
		url,
		config.AppConfig.MinioBucket,
		"survey",
		survey.Image,
	)

	return &survey, nil
}

func GetAllSurveys(ctx context.Context, categoryID int, name string) ([]model.SurveyResponse, error) {
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

	var surveys []model.SurveyResponse

	for rows.Next() {
		var survey model.SurveyResponse
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
		survey.Image = fmt.Sprintf("http://%s/%s/%s/%s",
			url,
			config.AppConfig.MinioBucket,
			"survey",
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

func GetAllCategories(ctx context.Context) ([]model.CategoryResponse, error) {
	query := `SELECT id, name FROM categories`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		log.Println("Error querying category:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []model.CategoryResponse

	for rows.Next() {
		var category model.CategoryResponse
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
