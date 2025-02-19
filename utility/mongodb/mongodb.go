package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	URI      string
	Database string
}

func NewMongoClient(ctx context.Context, cfg Config) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	// Ping the database
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
