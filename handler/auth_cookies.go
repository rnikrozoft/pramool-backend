package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const shortSessionRefreshSec = 3600 // 1 hour when "remember me" is off

// SessionCookieMaxAges returns access + refresh cookie Max-Age for login/refresh.
// Access stays at the configured hour; refresh is 1 hour or 7 days when remember is set.
func SessionCookieMaxAges(remember bool, accessDefaultSec, longRefreshDefaultSec int) (accessSec, refreshSec int) {
	if accessDefaultSec <= 0 {
		accessDefaultSec = shortSessionRefreshSec
	}
	if longRefreshDefaultSec <= 0 {
		longRefreshDefaultSec = shortSessionRefreshSec * 24 * 7
	}
	accessSec = accessDefaultSec
	if remember {
		refreshSec = longRefreshDefaultSec
	} else {
		refreshSec = shortSessionRefreshSec
	}
	return accessSec, refreshSec
}

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
