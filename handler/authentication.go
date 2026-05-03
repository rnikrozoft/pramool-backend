package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/service"
)

type AuthenticationHandler struct {
	validate              *validator.Validate
	authenticationService service.AuthenticationService
	accessCookieMaxAgeSec  int
	refreshCookieMaxAgeSec int
}

func NewAuthenticationHandler(
	validate *validator.Validate,
	authenticationService service.AuthenticationService,
	accessCookieMaxAgeSec int,
	refreshCookieMaxAgeSec int,
) AuthenticationHandler {
	if accessCookieMaxAgeSec <= 0 {
		accessCookieMaxAgeSec = 3600
	}
	if refreshCookieMaxAgeSec <= 0 {
		refreshCookieMaxAgeSec = 3600 * 24 * 7
	}
	return AuthenticationHandler{
		validate:               validate,
		authenticationService:   authenticationService,
		accessCookieMaxAgeSec:   accessCookieMaxAgeSec,
		refreshCookieMaxAgeSec: refreshCookieMaxAgeSec,
	}
}

func (h AuthenticationHandler) LoginByTel(c *fiber.Ctx) error {
	req := new(dto.RequestOTP)
	if err := validate(c, h.validate, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	tokens, err := h.authenticationService.LoginByTel(c.Context(), req.Tel)
	if err != nil {
		return responseCommonError(c, err)
	}

	ApplyAuthCookies(c, tokens.Access, tokens.Refresh, h.accessCookieMaxAgeSec, h.refreshCookieMaxAgeSec)
	return c.SendStatus(fiber.StatusOK)
}

// Refresh rotates access + refresh cookies using refresh_token (POST /auth/refresh).
func (h AuthenticationHandler) Refresh(c *fiber.Ctx) error {
	raw := c.Cookies("refresh_token")
	if raw == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "missing refresh token"})
	}
	userID, err := h.authenticationService.ValidateRefreshToken(raw)
	if err != nil {
		ClearAuthCookies(c)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid refresh token"})
	}
	access, err := h.authenticationService.GenerateAccessToken(userID)
	if err != nil {
		return responseCommonError(c, err)
	}
	refresh, err := h.authenticationService.GenerateRefreshToken(userID)
	if err != nil {
		return responseCommonError(c, err)
	}
	ApplyAuthCookies(c, access, refresh, h.accessCookieMaxAgeSec, h.refreshCookieMaxAgeSec)
	return c.SendStatus(fiber.StatusOK)
}

func (h AuthenticationHandler) Logout(c *fiber.Ctx) error {
	ClearAuthCookies(c)
	return c.SendStatus(fiber.StatusOK)
}
