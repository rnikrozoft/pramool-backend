package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/exception"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/service"
)

func (h PrivacyHandler) ListDataProcessors(c *fiber.Ctx) error {
	items, err := h.privacyService.ListDataProcessors(c.Context())
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(dto.DataProcessorListResponse{Items: items})
}

func (h PrivacyHandler) RecordCookieConsent(c *fiber.Ctx) error {
	req := new(dto.CookieConsentRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(exception.BadRequest(err))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(exception.BadRequest(err))
	}
	userID, _ := c.Locals("user_id").(string)
	ip, ua := RequestMetaFromFiber(c)
	if err := h.privacyService.RecordCookieConsent(c.Context(), userID, "", ip, ua, *req); err != nil {
		if errors.Is(err, service.ErrInvalidCookieVersion) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "เวอร์ชันนโยบายคุกกี้ไม่ตรงกับปัจจุบัน"})
		}
		return responseCommonError(c, err)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h PrivacyHandler) GetMarketingConsent(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	optIn, err := h.privacyService.GetMarketingOptIn(c.Context(), userID)
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(dto.MarketingConsentResponse{MarketingOptIn: optIn})
}

func (h PrivacyHandler) UpdateMarketingConsent(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	req := new(dto.MarketingConsentRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(exception.BadRequest(err))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(exception.BadRequest(err))
	}
	ip, ua := RequestMetaFromFiber(c)
	if err := h.privacyService.UpdateMarketingConsent(c.Context(), userID, "", ip, ua, *req); err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(dto.MarketingConsentResponse{MarketingOptIn: req.MarketingOptIn})
}

func (h PrivacyHandler) ListRetentionJobs(c *fiber.Ctx) error {
	resp, err := h.privacyService.ListRetentionJobs(c.Context())
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(resp)
}
