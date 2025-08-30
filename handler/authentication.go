package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool.in.th-backend/model/dto"
	"github.com/rnikrozoft/pramool.in.th-backend/service"
)

type AuthenticationHandler struct {
	validate              *validator.Validate
	authenticationService service.AuthenticationService
}

func NewAuthenticationHandler(
	validate *validator.Validate,
	authenticationService service.AuthenticationService,
) AuthenticationHandler {
	return AuthenticationHandler{
		validate:              validate,
		authenticationService: authenticationService,
	}
}

func (h AuthenticationHandler) LoginByTel(c *fiber.Ctx) error {
	req := new(dto.RequestOTP)
	if err := validate(c, h.validate, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	token, err := h.authenticationService.LoginByTel(c.Context(), req.Tel)
	if err != nil {
		return responseCommonError(c, err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   3600,
	})
	return c.SendStatus(fiber.StatusOK)
}
