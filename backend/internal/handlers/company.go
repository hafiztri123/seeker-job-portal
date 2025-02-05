package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hafiztri123/internal/core/domain"
	"github.com/hafiztri123/internal/core/ports"
	"github.com/hafiztri123/internal/core/services"
	"github.com/hafiztri123/internal/middleware"
)

type CompanyHandler struct {
	companyService *services.CompanyService
}

func NewCompanyHandler(companyService *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: companyService}
}

func (h *CompanyHandler) Register(c *fiber.Ctx) error {
    var req ports.CompanyRegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return &middleware.AppError{
            Code:    middleware.ErrInvalidInput,
            Message: "Invalid request format",
            Details: err.Error(),
        }
    }

    if err := h.companyService.Register(&req); err != nil {
        // Log the error for debugging
        fmt.Printf("Company registration error: %v\n", err)
        
        return &middleware.AppError{
            Code:    middleware.ErrInternalServer,
            Message: "Failed to register company",
            Details: err.Error(),
        }
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Company registered successfully",
    })
}

func (h *CompanyHandler) UpdateProfile (ctx *fiber.Ctx) error {
	company := ctx.Locals("company").(*domain.Company)
	var req ports.CompanyUpdateProfileRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	updatedCompany, err := h.companyService.UpdateProfile(company.Id.String(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(updatedCompany)
}

func (h *CompanyHandler) GetAllCompany (ctx *fiber.Ctx) error {
	companies, err := h.companyService.GetAllCompany()
	if err != nil {
		return err
	}

	return ctx.JSON(companies)
}

func (h *CompanyHandler) GetCompany (ctx *fiber.Ctx) error {
	company := ctx.Locals("company").(*domain.Company)
	return ctx.JSON(company)
}

