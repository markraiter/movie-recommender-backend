package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AccessCookieName  = "access-cookie"
	RefreshCookieName = "refresh-cookie"
)

type Claims struct {
	jwt.StandardClaims
	UserID primitive.ObjectID
}

type TokenPair struct {
	AccessToken   string `json:"access_token"`
	AccessExpire  time.Time
	RefreshToken  string `json:"refresh_token"`
	RefresgExpire time.Time
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email" example:"admin@example.com"`
	Password string `json:"password" validate:"required" example:"password12345"`
}
