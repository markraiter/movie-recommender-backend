package storage

import (
	"context"

	"github.com/markraiter/movie-recommender-backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	usersCollection = "users"
)

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

func New(ctx context.Context, cfg config.Mongo) (*MongoDB, error) {
	client, err := mongo.Connect(ctx,
		options.Client().
			ApplyURI(cfg.ConnectionString).SetAuth(
			options.Credential{
				Username: cfg.Username,
				Password: cfg.Password,
			},
		))
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.NameDB)

	return &MongoDB{client: client, db: db}, nil
}
