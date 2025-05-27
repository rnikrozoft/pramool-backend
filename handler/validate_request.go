package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool.in.th-backend/constant"
	"github.com/rnikrozoft/pramool.in.th-backend/exception"
)

func validate(c *fiber.Ctx, validate *validator.Validate, req interface{}) *exception.Response {
	if err := c.BodyParser(req); err != nil {
		e := exception.Set(constant.ErrorBadRequest)
		return &e
	}

	if err := validate.Struct(req); err != nil {
		e := exception.Set(constant.ErrorInvalidArguments)
		return &e
	}
	return nil
}
