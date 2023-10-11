package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/gofiber/swagger"
	"github.com/markraiter/movie-recommender-backend/config"
	"github.com/markraiter/movie-recommender-backend/internal/api/handlers"
)

const apiPrefixV1 = "/api"

func (s Server) initRoutes(app *fiber.App, cfg config.Config) {
	// identify := middlewares.NewUserIdentity(cfg)

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get(apiPrefixV1+"/health", handlers.HandlerAPIHealth())

	api := app.Group(apiPrefixV1)
	{
		api.Post("/sign-up", timeout.NewWithContext(handlers.HandlerSignUp(cfg, s.services, s.responser), cfg.Server.AppWriteTimeout))
		api.Post("/sign-in", timeout.NewWithContext(handlers.HandlerSignIn(cfg, s.services, s.responser), cfg.Server.AppWriteTimeout))
		api.Post("/refresh", timeout.NewWithContext(handlers.HandlerRefresh(cfg, s.services, *s.responser), cfg.Server.AppWriteTimeout))
	}
}
