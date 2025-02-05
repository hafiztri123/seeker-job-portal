package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/hafiztri123/internal/core/ports"
)


func AuthMiddleware(authService ports.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := extractToken(c)

		if tokenString == "" {
			return &AppError{
				Code: ErrUnauthorized,
				Message: "Missing token",
			}
		}

		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			return &AppError{
				Code: ErrUnauthorized,
				Message: "Invalid token",
			}
		}
		if claims["role"] == "company" {
			company, err := authService.GetCompany(claims["sub"].(string))
			if err != nil {
				return &AppError{
					Code: ErrUnauthorized,
					Message: "Invalid token",
				}
			}
			c.Locals("company", company)
		} else {
			user, err := authService.GetUser(claims["sub"].(string))
			if err != nil {
				return &AppError{
					Code: ErrUnauthorized,
					Message: "Invalid token",
				}
			}
			c.Locals("user", user)
		}

		return c.Next()

	}
}

func extractToken(c *fiber.Ctx) string {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return ""
    }
    
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return ""
    }
    
    return parts[1]
}