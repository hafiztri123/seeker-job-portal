package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	app.Get("/health", HealthCheck)

	t.Run("Returns200OK", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("ReturnCorrectPaylod", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, _ := app.Test(req)

		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "UP", result["status"])
		assert.NotNil(t, result["timestamp"])
	})
}