package models

import "github.com/google/uuid"

// Структуру пользователя
type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
}

// Данные из формы
type FormData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// Структура для ответа API
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

// DeleteRequest представляет запрос на удаление пользователей
type DeleteRequest struct {
	UserIDs []int `json:"user_ids"`
}

// DeleteResponse представляет ответ на удаление пользователей
type DeleteResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	DeletedCount int    `json:"deleted_count,omitempty"`
}
