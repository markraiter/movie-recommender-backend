package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/markraiter/movie-recommender-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoDB) GetEmail(email string) string {
	user := new(models.User)

	filter := bson.M{"email": email}

	err := m.db.Collection(usersCollection).FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return ""
	}

	return user.Email
}

func (m *MongoDB) Create(user *models.User) (primitive.ObjectID, error) {
	_, err := m.db.Collection(usersCollection).InsertOne(context.Background(), user)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("storage Create() error: %w", err)
	}

	filter := bson.M{"email": user.Email}
	result := m.db.Collection(usersCollection).FindOne(context.Background(), filter)
	if result.Err() != nil {
		return primitive.NilObjectID, fmt.Errorf("storage Create() error: %w", result.Err())
	}

	var createdUser models.User
	if err := result.Decode(&createdUser); err != nil {
		return primitive.NilObjectID, fmt.Errorf("storage Create() error: %w", err)
	}

	return createdUser.ID, nil
}

func (m *MongoDB) GetUserByEmail(email, password string) (*models.User, error) {
	var user models.User

	filter := bson.M{"email": email, "password": password}
	result := m.db.Collection(usersCollection).FindOne(context.Background(), filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, models.ErrWrongCredentials
		}

		return nil, fmt.Errorf("storage GetUserByEmail() error: %w", result.Err())
	}

	if err := result.Decode(&user); err != nil {
		return nil, fmt.Errorf("storage GetUserByEmail() error: %w", err)
	}

	return &user, nil
}
