package model

type CreateSurveyRequest struct {
	Name        string `form:"name" binding:"required"`
	Price       int    `form:"price" binding:"required"`
	Description string `form:"description" binding:"required"`
	CategoryID  int    `form:"category_id" binding:"required"`
}

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

type SurveyEvent struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Image       string `json:"image"`
	CategoryID  int    `json:"category_id"`
}
