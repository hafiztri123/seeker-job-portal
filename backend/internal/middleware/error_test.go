package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	app := fiber.New(
		fiber.Config{
			ErrorHandler: ErrorHandler,
		},
	)

	app.Get("/error-test", func (c *fiber.Ctx) error  {
		return fiber.NewError(fiber.StatusBadRequest, "test error")
	})

	app.Get("/custom-error", func (c *fiber.Ctx) error  {
		return &AppError{
			Code: ErrInvalidInput,
			Message: "validation failed",
			Details: map[string]string{"field": "invalid"},
		}
	})

	t.Run("HandlesFiberError", func (t *testing.T)  {
		req := httptest.NewRequest("GET", "/error-test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "test error", result["message"])
	})

	t.Run("HandlesCustomError", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/custom-error", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "INVALID_INPUT", result["code"])
		assert.Equal(t, "validation failed", result["message"])
	})
}