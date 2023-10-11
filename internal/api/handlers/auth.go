package handlers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/markraiter/movie-recommender-backend/config"
	"github.com/markraiter/movie-recommender-backend/internal/api/middlewares"
	"github.com/markraiter/movie-recommender-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type authServices interface {
	GetEmail(email string) string
	CreateUser(cfg config.Config, user *models.User) (primitive.ObjectID, error)
	GetUserByEmail(email, password string) (*models.User, error)
	GetTokenPair(cfg config.Config, email, password string) (*models.TokenPair, error)
	Refresh(userID primitive.ObjectID, cfg config.Config) (*models.TokenPair, error)
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body models.User true "account info"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 406 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /sign-up [post].
func HandlerSignUp(cfg config.Config, s authServices, resp *models.Responser) fiber.Handler {
	return func(c *fiber.Ctx) error {
		message := "HandlerSignUp"
		user := new(models.User)

		if err := c.BodyParser(&user); err != nil {
			return resp.Debug(message, fiber.StatusBadRequest, err)
		}

		if err := validate.Struct(user); err != nil {
			return resp.Debug(message, fiber.StatusNotAcceptable, err)
		}

		id, err := s.CreateUser(cfg, user)
		if err != nil {
			if errors.Is(err, models.ErrUniqueViolation) {
				return resp.Warn(message, fiber.StatusNotAcceptable, err)
			}

			return resp.Error(message, fiber.StatusNotAcceptable, err)
		}

		return c.Status(fiber.StatusCreated).JSON(models.Response{Message: "user successfully created with ID: " + id.Hex()})
	}
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body models.LoginInput true "credentials"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /sign-in [post].
func HandlerSignIn(cfg config.Config, s authServices, resp *models.Responser) fiber.Handler {
	return func(c *fiber.Ctx) error {
		message := "HandlerSignIn"
		var input models.LoginInput

		if err := c.BodyParser(&input); err != nil {
			return resp.Debug(message, fiber.StatusNotAcceptable, err)
		}

		tokenPair, err := s.GetTokenPair(cfg, input.Email, input.Password)
		if err != nil {
			if errors.Is(err, models.ErrWrongCredentials) {
				return resp.Warn(message, fiber.StatusUnauthorized, err)
			}

			return resp.Error(message, fiber.StatusNotAcceptable, err)
		}

		accessCookie := newCookie(models.AccessCookieName, tokenPair.AccessToken, tokenPair.AccessExpire)

		c.Cookie(accessCookie)

		refreshCookie := newCookie(models.RefreshCookieName, tokenPair.RefreshToken, tokenPair.RefresgExpire)
		refreshCookie.Path = "/api/v1/refresh" // NEED DOUBLE CHECK

		c.Cookie(refreshCookie)

		return c.Status(fiber.StatusOK).JSON(models.Response{Message: "you are logged in"})
	}
}

// @Summary Refresh
// @Tags auth
// @Description refresh
// @ID refresh
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /refresh [post].
func HandlerRefresh(cfg config.Config, s authServices, resp models.Responser) fiber.Handler {
	return func(c *fiber.Ctx) error {
		message := "HandlerRefresh"

		refreshString := c.Cookies(models.RefreshCookieName)

		if refreshString == "" {
			return resp.Warn(message, fiber.StatusForbidden, models.ErrEmptyCookie)
		}

		userID, err := middlewares.ParseToken(refreshString, cfg.Auth.SigningKey)
		if err != nil {
			return resp.Debug(message, fiber.StatusBadRequest, err)
		}

		tokenPair, err := s.Refresh(userID, cfg)
		if err != nil {
			if errors.Is(err, models.ErrNoSession) {
				return resp.Warn(message, fiber.StatusUnauthorized, err)
			}

			return resp.Error(message, fiber.StatusInternalServerError, err)
		}

		accessCookie := newCookie(
			models.AccessCookieName,
			tokenPair.AccessToken,
			tokenPair.AccessExpire,
		)

		refreshCookie := newCookie(
			models.RefreshCookieName,
			tokenPair.RefreshToken,
			tokenPair.RefresgExpire,
		)

		refreshCookie.Path = "/api/v1/refresh" // NEED DOUBLE CHECK

		c.Cookie(accessCookie)
		c.Cookie(refreshCookie)

		return c.Status(fiber.StatusOK).JSON(models.Response{Message: "refreshed"})
	}
}

func newCookie(name, value string, expire time.Time) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expire,
		Secure:   false,
		HTTPOnly: true,
		SameSite: "none",
	}
}
