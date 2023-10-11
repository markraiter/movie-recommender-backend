package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/markraiter/movie-recommender-backend/models"
)

var validate = validator.New() //nolint:gochecknoglobals

// @Summary Show the status of server.
// @Description Ping health of API for Docker.
// @Tags Health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get].
func HandlerAPIHealth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(models.Response{Message: "healthy"})
	}
}
