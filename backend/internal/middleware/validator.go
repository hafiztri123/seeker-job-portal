package middleware

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type ValidationError struct {
	Field string `json:"field"`
	Tag string `json:"tag"`
	Value string `json:"value"`
}


func init() {
	validate.RegisterValidation("password_strength", validatePasswordStrength)
}

func validatePasswordStrength(fl validator.FieldLevel) bool {
	password :=fl.Field().String()
	return len(password) >= 8 &&
		strings.ContainsAny(password, "0123456789") &&
		strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}


func ValidateRequest(payload interface{}) []*ValidationError {
	var errors []*ValidationError
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}

func ValidateBody(payload interface{}) fiber.Handler {
	return func (c *fiber.Ctx) error {
		if err := c.BodyParser(payload); err != nil {
			return &AppError{
				Code: ErrInvalidInput,
				Message: "Invalid request body",
			}
		}

		if errors := ValidateRequest(payload); len(errors) > 0 {
			return &AppError{
				Code: ErrInvalidInput,
				Message: "validation failed",
				Details: errors,
			}
		}

		c.Locals("validated", payload)
		return c.Next()
		
	}
}