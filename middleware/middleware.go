package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rnikrozoft/pramool.in.th-backend/constant"
	"github.com/rnikrozoft/pramool.in.th-backend/model"
)

type Middleware struct {
	JwtSecret string
}

func (m Middleware) JWTMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("access_token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing token")
	}

	token, err := jwt.ParseWithClaims(cookie, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}

	claims := token.Claims.(*model.CustomClaims)
	c.Locals(constant.ClaimsJWTObject, claims.UserID)

	return c.Next()
}
