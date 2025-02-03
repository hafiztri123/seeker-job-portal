package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hafiztri123/internal/core/domain"
	"github.com/hafiztri123/internal/core/ports"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type loginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

type RegisterRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
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