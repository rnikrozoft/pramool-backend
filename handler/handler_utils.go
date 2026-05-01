package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/exception"
)

func responseCommonError(c *fiber.Ctx, err error) error {
	var appErr *exception.AppError
	if errors.As(err, &appErr) {
		return c.Status(appErr.StatusCode).JSON(appErr)
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": err.Error(),
	})
}

func validate(c *fiber.Ctx, validate *validator.Validate, req interface{}) error {
	if err := c.BodyParser(req); err != nil {
		return exception.BadRequest(err)
	}

	if err := validate.Struct(req); err != nil {
		return exception.BadRequest(err)
	}
	return nil
}
