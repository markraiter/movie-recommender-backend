package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id" validate:"omitempty" example:""`
	Email     string             `bson:"email" json:"email" validate:"email" example:"admin@example.com"`
	Password  string             `bson:"password" json:"password" validate:"min=8" example:"password12345"`
	FirstName string             `bson:"first_name" json:"first_name" validate:"omitempty" example:"Alex"`
	LastName  string             `bson:"last_name" json:"last_name" validate:"omitempty" example:"Johnson"`
}
