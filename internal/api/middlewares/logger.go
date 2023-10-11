package middlewares

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func NewLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		chainErr := c.Next()

		code := fiber.StatusInternalServerError
		if chainErr != nil {
			var localErr *fiber.Error
			if errors.As(chainErr, &localErr) {
				code = localErr.Code
			}

			c.Status(code)
		}

		logger.Info("access_log", "metod", c.Method(), "status", c.Response().StatusCode(), "path", c.OriginalURL())

		return chainErr
	}
}
