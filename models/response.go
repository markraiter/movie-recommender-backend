package models

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Message string `json:"message" example:"response message"`
}

type Responser struct {
	Log *slog.Logger
}

func (r *Responser) Error(msg string, status int, err error) *fiber.Error {
	errString := fmt.Sprintf("%s error: %s", msg, err.Error())

	r.Log.Error("error_log", "message", errString)

	return fiber.NewError(status, errString)
}

func (r *Responser) Warn(msg string, status int, err error) *fiber.Error {
	errString := fmt.Sprintf("%s error: %s", msg, err.Error())

	r.Log.Warn("error_log", "message", errString)

	return fiber.NewError(status, errString)
}

func (r *Responser) Info(msg string, status int, err error) *fiber.Error {
	errString := fmt.Sprintf("%s error: %s", msg, err.Error())

	r.Log.Info("error_log", "message", errString)

	return fiber.NewError(status, errString)
}

func (r *Responser) Debug(msg string, status int, err error) *fiber.Error {
	errString := fmt.Sprintf("%s error: %s", msg, err.Error())

	r.Log.Debug("error_log", "message", errString)

	return fiber.NewError(status, errString)
}
