package user

import (
	"context"
	"log"
	"user/pkg/database"

	"github.com/jackc/pgx/v5"
)

type UserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func StoreUser(ctx context.Context, req UserInput) error {
	query := `INSERT INTO users (name, email, password) 
			  VALUES (@name, @email, @password)`
	args := pgx.NamedArgs{
		"name":     req.Name,
		"email":    req.Email,
		"password": req.Password,
	}

	_, err := database.DB.Exec(ctx, query, args)
	if err != nil {
		log.Println("Error inserting user:", err)
		return err
	}

	return nil
}

func GetUserByID(ctx context.Context, id int) (*UserResponse, error) {
	query := `SELECT id, name, email, password FROM users WHERE id = $1`

	row := database.DB.QueryRow(ctx, query, id)

	var user UserResponse
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Println("Error fetching user:", err)
		return nil, err
	}

	return &user, nil
}

func GetAllUsers(ctx context.Context) ([]UserResponse, error) {
	query := `SELECT id, name, email, password FROM users`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		log.Println("Error querying users:", err)
		return nil, err
	}
	defer rows.Close()

	var users []UserResponse

	for rows.Next() {
		var user UserResponse
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			log.Println("Error scanning user:", err)
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Println("Rows iteration error:", err)
		return nil, err
	}

	return users, nil
}
