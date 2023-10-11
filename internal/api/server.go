package api

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/markraiter/movie-recommender-backend/config"
	"github.com/markraiter/movie-recommender-backend/internal/api/middlewares"
	"github.com/markraiter/movie-recommender-backend/internal/services"
	"github.com/markraiter/movie-recommender-backend/models"
)

type Server struct {
	HTTPServer *fiber.App
	services   *services.Services
	responser  *models.Responser
}

func NewServer(cfg config.Config, s *services.Services, log *slog.Logger) *Server {
	server := new(Server)
	server.responser = &models.Responser{Log: log}
	server.services = s

	fconfig := fiber.Config{
		ReadTimeout:  cfg.Server.AppReadTimeout,
		WriteTimeout: cfg.Server.AppWriteTimeout,
		IdleTimeout:  cfg.Server.AppIdleTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var localError *fiber.Error
			if errors.As(err, &localError) {
				code = localError.Code
			}

			c.Status(code)

			if err := c.JSON(models.Response{Message: localError.Message}); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			return nil
		},
	}
	server.HTTPServer = fiber.New(fconfig)

	server.HTTPServer.Use(cors.New(corsConfig()))

	server.HTTPServer.Use(recover.New())

	server.HTTPServer.Use(middlewares.NewLogger(log))

	server.initRoutes(server.HTTPServer, cfg)

	return server
}

func corsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     `*`,
		AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Credentials",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowCredentials: true,
	}
}
