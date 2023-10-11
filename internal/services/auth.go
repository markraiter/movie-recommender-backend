package services

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/markraiter/movie-recommender-backend/config"
	"github.com/markraiter/movie-recommender-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthInterface interface {
	GetEmail(email string) string
	Create(user *models.User) (primitive.ObjectID, error)
	GetUserByEmail(email, password string) (*models.User, error)
}

type Auth struct {
	DB AuthInterface
}

func (a *Auth) GetTokenPair(cfg config.Config, email, password string) (*models.TokenPair, error) {
	user, err := a.DB.GetUserByEmail(email, generatePasswordHash(cfg, password))
	if err != nil {
		return nil, fmt.Errorf("service GetTokenPair() error: %w", err)
	}

	accessExpire := time.Now().Add(cfg.Auth.AccessTokenTTL)
	refreshExpire := time.Now().Add(cfg.Auth.RefreshTokenTTL)

	accessToken, err := generateJWT(accessExpire, cfg.Auth.SigningKey, user.ID)
	if err != nil {
		return nil, fmt.Errorf("service GetTokenPair() error: %w", err)
	}

	refreshToken, err := generateJWT(refreshExpire, cfg.Auth.SigningKey, user.ID)
	if err != nil {
		return nil, fmt.Errorf("service GetTokenPair() error: %w", err)
	}

	tokenPair := models.TokenPair{
		AccessToken:   accessToken,
		AccessExpire:  accessExpire,
		RefreshToken:  refreshToken,
		RefresgExpire: refreshExpire,
	}

	return &tokenPair, nil
}

func generateJWT(expire time.Time, signingKey string, userID primitive.ObjectID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: userID,
	})

	return token.SignedString([]byte(signingKey))
}

func (a *Auth) Refresh(userID primitive.ObjectID, cfg config.Config) (*models.TokenPair, error) {
	accessExpire := time.Now().Add(cfg.Auth.AccessTokenTTL)
	refreshExpire := time.Now().Add(cfg.Auth.RefreshTokenTTL)

	accessToken, err := generateJWT(accessExpire, cfg.Auth.SigningKey, userID)
	if err != nil {
		return nil, fmt.Errorf("service Refresh() error: %w", err)
	}

	refreshToken, err := generateJWT(refreshExpire, cfg.Auth.SigningKey, userID)
	if err != nil {
		return nil, fmt.Errorf("service Refresh() error: %w", err)
	}

	tokenPair := models.TokenPair{
		AccessToken:   accessToken,
		AccessExpire:  accessExpire,
		RefreshToken:  refreshToken,
		RefresgExpire: refreshExpire,
	}

	return &tokenPair, nil
}

func (a *Auth) GetEmail(email string) string {
	res := a.DB.GetEmail(email)

	return res
}

func (a *Auth) CreateUser(cfg config.Config, user *models.User) (primitive.ObjectID, error) {
	email := a.GetEmail(user.Email)

	if user.Email == email {
		return primitive.NilObjectID, fmt.Errorf("this user is already exists: %w", models.ErrUniqueViolation)
	}

	user.Password = generatePasswordHash(cfg, user.Password)

	id, err := a.DB.Create(user)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("service CreateUser() error: %w", err)
	}

	return id, nil
}

func (a *Auth) GetUserByEmail(email, password string) (*models.User, error) {
	usr, err := a.DB.GetUserByEmail(email, password)
	if err != nil {
		return nil, fmt.Errorf("service GetUserByEmail() error: %w", err)
	}

	return usr, nil
}

func generatePasswordHash(cfg config.Config, password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(cfg.Auth.Salt)))
}
