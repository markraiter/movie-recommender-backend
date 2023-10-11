package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/markraiter/movie-recommender-backend/config"
	_ "github.com/markraiter/movie-recommender-backend/docs"
	"github.com/markraiter/movie-recommender-backend/internal/api"
	"github.com/markraiter/movie-recommender-backend/internal/services"
	"github.com/markraiter/movie-recommender-backend/internal/storage"
)

//	@title			MOVIE-RECOMMENDER API
//	@version		1.0
//	@description	Docs for movie-recommender backend API
//	@contact.name	Mark Raiter
//	@contact.email	raitermrk@gmail.com
//  @host  			localhost:8000
//	@BasePath		/api

func main() {
	const timoutLimit = 5

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.InitConfig()
	if err != nil {
		log.Error("InitConfig", "error", err.Error())

		return
	}

	db, err := storage.New(context.Background(), cfg.DB)
	if err != nil {
		log.Error("New storage", "error", err.Error())

		return
	}

	service := services.NewService(db)

	server := api.NewServer(cfg, &service, log)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		if err := server.HTTPServer.Listen(cfg.Server.AppAddress); err != nil {
			log.Error("HTTPServer.Listen", "error", err.Error())
		}
	}()

	<-quit

	if err := server.HTTPServer.ShutdownWithTimeout(timoutLimit * time.Second); err != nil {
		log.Error("ShutdownWithTimeout", "error", err.Error())
	}

	if err := server.HTTPServer.Shutdown(); err != nil {
		log.Error("Shutdown", "error", err.Error())
	}

	log.Info("server stopped")
}
