package database

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func CreateClient(ctx context.Context, uri string) (*mongo.Database, error) {
	log.Info("connecting to mongodb")

	client, error := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAppName("klein"))

	if error != nil {
		log.Fatal(error)
	}

	if error := client.Ping(ctx, readpref.Primary()); error != nil {
		log.Fatal("couldn't establish a connection to mongodb")
		return nil, error
	}

	log.Info("successfully connected to mongodb")
	return client.Database("klein"), nil
}
