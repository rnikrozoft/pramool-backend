package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// ApplyAuthCookies sets HttpOnly access_token + refresh_token (same-site lax, path /).
func ApplyAuthCookies(c *fiber.Ctx, accessToken, refreshToken string, accessMaxAgeSec, refreshMaxAgeSec int) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   accessMaxAgeSec,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   refreshMaxAgeSec,
	})
}

// ClearAuthCookies clears access + refresh cookies.
func ClearAuthCookies(c *fiber.Ctx) {
	exp := time.Unix(0, 0)
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
		Expires:  exp,
		MaxAge:   -1,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
		Expires:  exp,
		MaxAge:   -1,
	})
}
