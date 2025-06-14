package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool.in.th-backend/model"
	"github.com/rnikrozoft/pramool.in.th-backend/model/dto"
	"github.com/rnikrozoft/pramool.in.th-backend/service"
)

type OtpHandler struct {
	validate        *validator.Validate
	otpService      service.OTPService
	registerService service.RegisterService
}

func NewOTPHandler(
	validate *validator.Validate,
	otpService service.OTPService,
	registerService service.RegisterService,
) OtpHandler {
	return OtpHandler{
		validate:        validate,
		otpService:      otpService,
		registerService: registerService,
	}
}

func (h OtpHandler) RequestOTP(c *fiber.Ctx) error {
	req := new(dto.RequestOTP)
	if err := validate(c, h.validate, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	// if err := h.registerService.RegisterTelIfNotExist(c.Context(), req.Tel); err != nil {
	// 	return responseCommonError(c, err)
	// }

	// res, err := h.otpService.Request(req.Tel)
	// if err != nil {
	// 	return responseCommonError(c, err)
	// }
	// return c.Status(fiber.StatusOK).JSON(res)

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"status": "success",
		"token":  "36nmr8Jbw9vz5aeuXIv5olLX4xBZeg2j",
		"refNo":  "YZJRN",
	})
}

func (h OtpHandler) VerifyOTP(c *fiber.Ctx) error {
	req := new(model.VerifyRequest)
	if err := validate(c, h.validate, req); err != nil {
		return responseCommonError(c, err)
	}

	// if err := h.otpService.Verify(req.Token, req.PIN); err != nil {
	// 	return responseCommonError(c, err)
	// }
	return c.SendStatus(200)
}
