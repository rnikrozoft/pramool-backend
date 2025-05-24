package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool.in.th-backend/exception"
	"github.com/rnikrozoft/pramool.in.th-backend/model"
	"github.com/rnikrozoft/pramool.in.th-backend/model/dto"
	"github.com/rnikrozoft/pramool.in.th-backend/service"
)

type RegisterHandler struct {
	validate        *validator.Validate
	registerService service.RegisterService
}

func NewRegisterHandler(validate *validator.Validate, registerService service.RegisterService) RegisterHandler {
	return RegisterHandler{
		validate:        validate,
		registerService: registerService,
	}
}

// Register godoc
// @Summary      Register
// @Description  User Register
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      dto.User  true  "User data"
// @Success      201   {object}  model.RegisterResponse
// @Router       /register [post]
func (h RegisterHandler) Register(c *fiber.Ctx) error {
	user := new(dto.User)
	if err := c.BodyParser(user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if err := h.validate.Struct(user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	token, err := h.registerService.Register(c.Context(), *user)
	if err != nil {
		e := exception.Set(fiber.ErrBadRequest.Error(), err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(e)
	}
	return c.Status(fiber.StatusCreated).JSON(model.RegisterResponse{
		Token: token,
	})
}
