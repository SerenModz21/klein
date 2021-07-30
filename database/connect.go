package database

import (
	"context"

	"github.com/mediocregopher/radix/v4"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongo(ctx context.Context, uri string) (*mongo.Database, error) {
	log.Info("connecting to mongoDB")

	client, error := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAppName("klein"))

	if error != nil {
		log.Fatal(error)
	}

	if error := client.Ping(ctx, readpref.Primary()); error != nil {
		log.Fatal("couldn't establish a connection to mongoDB")
		return nil, error
	}

	log.Info("successfully connected to mongoDB")
	return client.Database("klein"), nil
}

func ConnectRedis(ctx context.Context, uri string) (radix.Client, error) {
	log.Info("connecting to redis")

	client, error := (radix.PoolConfig{}).New(ctx, "tcp", uri)
	if error != nil {
		log.Fatal(error)
	}

	log.Info("successfully connected to redis")
	return client, nil
}
