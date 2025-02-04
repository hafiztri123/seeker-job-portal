package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hafiztri123/internal/core/domain"
	"github.com/hafiztri123/internal/core/ports"
	"github.com/hafiztri123/internal/middleware"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}


	tokens, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(middleware.AppError{
			Code: middleware.ErrInvalidCredentials,
			Message: "Invalid credentials",
			
		})
	}

	return c.JSON(loginResponse{
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

type RegisterRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password_strength"`
	FullName string `json:"full_name" validate:"required"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	user, err := domain.NewUser(req.Email, req.FullName, req.Password)
	if err != nil {
		return err
	}

	if err := h.authService.Register(user); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})

}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if errors := middleware.ValidateRequest(&req); len(errors) > 0 {
		return &middleware.AppError{
			Code: middleware.ErrInvalidInput,
			Message: "validation failed",
			Details: errors,
		}
	}

	tokens, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return err
	}

	return c.JSON(loginResponse{
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

