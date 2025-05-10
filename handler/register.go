package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool.in.th-backend/model"
	"github.com/rnikrozoft/pramool.in.th-backend/service"
)

type RegisterHandler struct {
	validate        *validator.Validate
	registerService service.Register
}

func NewRegisterHandler(validate *validator.Validate, registerService service.Register) RegisterHandler {
	return RegisterHandler{
		validate:        validate,
		registerService: registerService,
	}
}

func (h RegisterHandler) Register(c *fiber.Ctx) error {
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if err := h.validate.Struct(user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.registerService.Register(*user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": err,
		})
	}
	return c.SendStatus(fiber.StatusCreated)
}
