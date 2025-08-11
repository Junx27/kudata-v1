package model

type CreateSurveyRequest struct {
	Name        string `form:"name" binding:"required"`
	Price       int    `form:"price" binding:"required"`
	Description string `form:"description" binding:"required"`
	CategoryID  int    `form:"category_id" binding:"required"`
}
type MessageSurvey struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Image       string `json:"image"`
	CategoryID  int    `json:"category_id"`
}
