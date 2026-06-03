package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/service"
)

type PasswordResetHandler struct {
	validate             *validator.Validate
	passwordResetService service.PasswordResetService
}

func NewPasswordResetHandler(validate *validator.Validate, passwordResetService service.PasswordResetService) PasswordResetHandler {
	return PasswordResetHandler{
		validate:             validate,
		passwordResetService: passwordResetService,
	}
}

func (h PasswordResetHandler) Check(c *fiber.Ctx) error {
	req := new(dto.ForgotPasswordCheckRequest)
	if err := validate(c, h.validate, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	status, tel, err := h.passwordResetService.CheckEligibility(c.Context(), req.Tel)
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(fiber.Map{
		"status": status,
		"tel":    tel,
	})
}

func (h PasswordResetHandler) Reset(c *fiber.Ctx) error {
	req := new(dto.ForgotPasswordResetRequest)
	if err := validate(c, h.validate, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	if err := h.passwordResetService.ResetPassword(c.Context(), req.Tel, req.Token, req.PIN, req.Password); err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(fiber.Map{"message": "ตั้งรหัสผ่านใหม่สำเร็จ"})
}
