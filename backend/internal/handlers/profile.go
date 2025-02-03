package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hafiztri123/internal/core/domain"
	"github.com/hafiztri123/internal/core/ports"
)

type ProfileHandler struct {
	userService ports.UserService
}


func NewProfileHandler(userService ports.UserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	return c.JSON(user)
}

func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	var req ports.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	user := c.Locals("user").(*domain.User)

	updatedUser, err := h.userService.UpdateProfile(user.ID, req)
	if err != nil {
		return err
	}

	return c.JSON(updatedUser)
}

