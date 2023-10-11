package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/markraiter/movie-recommender-backend/config"
	"github.com/markraiter/movie-recommender-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const AccessCookieName = "access-cookie"

func NewUserIdentity(cfg config.Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Cookies(AccessCookieName)

		if tokenString == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "empty cookie")
		}

		userID, err := ParseToken(tokenString, cfg.Auth.SigningKey)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		ctx.Locals("userID", userID)
		ctx.Locals("refreshString", tokenString)

		return ctx.Next()
	}
}

func ParseToken(tokenString, signingKey string) (primitive.ObjectID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.ErrInvalidSigningMethod
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("accessToken throws an error during parsing: %w", err)
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return primitive.NilObjectID, models.ErrWrongTokenClaimType
	}

	return claims.UserID, nil
}
