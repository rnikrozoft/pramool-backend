package exception

import (
	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(message string, statusCode int) *AppError {
	return &AppError{Message: message, StatusCode: statusCode}
}

func NotFound(err error) *AppError {
	return New(err.Error(), fiber.StatusNotFound)
}

func Internal(err error) *AppError {
	return New(err.Error(), fiber.StatusInternalServerError)
}

func BadRequest(err error) *AppError {
	return New(err.Error(), fiber.StatusBadRequest)
}

func Forbidden(message string) *AppError {
	return New(message, fiber.StatusForbidden)
}
