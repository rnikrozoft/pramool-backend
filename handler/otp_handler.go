package handler

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/model"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/service"
)

type OtpHandler struct {
	validate        *validator.Validate
	otpService      service.OTPService
	registerService service.RegisterService
	userService     service.UserService
}

func NewOTPHandler(
	validate *validator.Validate,
	otpService service.OTPService,
	registerService service.RegisterService,
	userService service.UserService,
) OtpHandler {
	return OtpHandler{
		validate:        validate,
		otpService:      otpService,
		registerService: registerService,
		userService:     userService,
	}
}

func (h OtpHandler) RequestOTP(c *fiber.Ctx) error {
	req := new(dto.RequestOTP)
	if err := validate(c, h.validate, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := h.registerService.RegisterTelIfNotExist(c.Context(), req.Tel); err != nil {
		return responseCommonError(c, err)
	}

	bannedUntil, err := h.userService.GetActiveOTPBanUntil(c.Context(), req.Tel)
	if err != nil {
		return responseCommonError(c, err)
	}
	if bannedUntil != nil {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"message":      "tel is temporarily banned",
			"banned_until": bannedUntil.Unix(),
		})
	}

	// res, err := h.otpService.Request(req.Tel)
	// if err != nil {
	// 	return responseCommonError(c, err)
	// }
	// return c.Status(fiber.StatusOK).JSON(res)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"token":   "36nmr8Jbw9vz5aeuXIv5olLX4xBZeg2j",
		"refNo":   "YZJRN",
		"channel": "tel",
	})
}

func (h OtpHandler) RecordTimeout(c *fiber.Ctx) error {
	req := new(dto.OTPTimeoutRequest)
	if err := validate(c, h.validate, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := h.registerService.RegisterTelIfNotExist(c.Context(), req.Tel); err != nil {
		return responseCommonError(c, err)
	}

	bannedUntil, timeoutCount, err := h.userService.RecordOTPTimeout(c.Context(), req.Tel)
	if err != nil {
		return responseCommonError(c, err)
	}

	if bannedUntil == nil {
		return c.JSON(fiber.Map{
			"status":        "ok",
			"timeout_count": timeoutCount,
			"channel":       "tel",
		})
	}

	return c.JSON(fiber.Map{
		"status":        "banned",
		"timeout_count": timeoutCount,
		"banned_until":  bannedUntil.Unix(),
		"banned_in_sec": int(time.Until(*bannedUntil).Seconds()),
		"channel":       "tel",
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
