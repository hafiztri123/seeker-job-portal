package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/hafiztri123/internal/core/ports"
)


func AuthMiddleware(authService ports.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return &AppError{
				Code: ErrUnauthorized,
				Message: "Missing autorization header",
			}
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return &AppError{
				Code: ErrUnauthorized,
				Message: "Invalid authorization header format",
			}
		}

		user, err := authService.ValidateToken(tokenParts[1])
		if err != nil {
			return &AppError{
				Code: ErrUnauthorized,
				Message: "Invalid token",
			}
		}

		c.Locals("user", user)
		return c.Next()
	}
}