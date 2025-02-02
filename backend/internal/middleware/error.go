package middleware

import "github.com/gofiber/fiber/v2"

type ErrorCode string

const (
	ErrInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrNotFound         ErrorCode = "NOT_FOUND"
	ErrUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrInternalServer   ErrorCode = "INTERNAL_SERVER"
)

type AppError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*AppError); ok {
		switch e.Code {
		case ErrInvalidInput:
			code = fiber.StatusBadRequest
		case ErrNotFound:
			code = fiber.StatusNotFound
		case ErrUnauthorized:
			code = fiber.StatusUnauthorized
		case ErrInternalServer:
			code = fiber.StatusInternalServerError
		}
		return c.Status(code).JSON(e)
	}

	if e, ok := err.(*fiber.Error); ok {
		return c.Status(e.Code).JSON(fiber.Map{
			"message": e.Message,
		})
	}

	return c.Status(code).JSON(fiber.Map{
		"message": "Internal Server Error",
	})
}