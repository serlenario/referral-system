// internal/models/error_response.go
package models

// ErrorResponse представляет структуру ошибки
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request parameters"`
}
