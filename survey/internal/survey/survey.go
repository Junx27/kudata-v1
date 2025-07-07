package survey

import (
	"context"
	"log"
	"survey/pkg/database"

	"github.com/jackc/pgx/v5"
)

type SurveyInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SurveyResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func StoreSurvey(ctx context.Context, req SurveyInput) error {
	query := `INSERT INTO surveys (name, description) 
			  VALUES (@name, @description)`
	args := pgx.NamedArgs{
		"name":        req.Name,
		"description": req.Description,
	}

	_, err := database.DB.Exec(ctx, query, args)
	if err != nil {
		log.Println("Error inserting survey:", err)
		return err
	}

	return nil
}

func GetSurveyByID(ctx context.Context, id int) (*SurveyResponse, error) {
	query := `SELECT id, name, description FROM surveys WHERE id = $1`

	row := database.DB.QueryRow(ctx, query, id)

	var survey SurveyResponse
	err := row.Scan(&survey.ID, &survey.Name, &survey.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Println("Error fetching survey:", err)
		return nil, err
	}

	return &survey, nil
}

func GetAllSurveys(ctx context.Context) ([]SurveyResponse, error) {
	query := `SELECT id, name, description FROM surveys`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		log.Println("Error querying survey:", err)
		return nil, err
	}
	defer rows.Close()

	var surveys []SurveyResponse

	for rows.Next() {
		var survey SurveyResponse
		err := rows.Scan(&survey.ID, &survey.Name, &survey.Description)
		if err != nil {
			log.Println("Error scanning survey:", err)
			return nil, err
		}
		surveys = append(surveys, survey)
	}
	if err = rows.Err(); err != nil {
		log.Println("Rows iteration error:", err)
		return nil, err
	}

	return surveys, nil
}
