package services

import (
	"context"
	"errors"

	"github.com/mediocregopher/radix/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sach.demiboy.me/database/models"
)

type IUrl interface {
	Insert(models.Url) (models.Url, error)
	Get(string) (models.Url, error)
	Delete(string, string) (models.Url, error)
	Search(interface{}, options.FindOptions) ([]models.Url, error)
}

type UrlClient struct {
	Ctx        context.Context
	Collection *mongo.Collection
	Redis      radix.Client
}

func (client *UrlClient) Insert(document models.Url) (models.Url, error) {
	url := models.Url{}

	_, error := client.Collection.InsertOne(client.Ctx, document)
	client.Redis.Do(client.Ctx, radix.Cmd(nil, "HSET", document.Slug, "DeletionKey", document.DeletionKey, "Long", document.Long, "Slug", document.Slug))

	if error != nil {
		return url, error
	}

	return client.Get(document.Slug)
}

func (client *UrlClient) Get(slug string) (models.Url, error) {
	url := models.Url{}
	cached := false

	client.Redis.Do(client.Ctx, radix.Cmd(&cached, "EXISTS", slug))

	if cached {
		client.Redis.Do(client.Ctx, radix.Cmd(&url, "HGETALL", slug))
		return url, nil
	} else {
		error := client.Collection.FindOne(client.Ctx, bson.M{"slug": slug}).Decode(&url)
		client.Redis.Do(client.Ctx, radix.Cmd(nil, "HSET", url.Slug, "DeletionKey", url.DeletionKey, "Long", url.Long, "Slug", url.Slug))

		return url, error
	}
}

func (client *UrlClient) Search(filter interface{}, options options.FindOptions) ([]models.Url, error) {
	urls := []models.Url{}

	if filter == nil {
		filter = bson.M{}
	}

	if options.Limit == nil {
		options.Limit = new(int64)
	}

	cursor, err := client.Collection.Find(client.Ctx, filter, &options)

	if err != nil {
		return urls, err
	}

	for cursor.Next(client.Ctx) {
		row := models.Url{}
		cursor.Decode(&row)
		urls = append(urls, row)
	}

	return urls, nil
}

func (client *UrlClient) Delete(slug string, deletionKey string) (models.Url, error) {
	url := models.Url{}

	if data, error := client.Get(slug); error != nil {
		return url, error
	} else {
		if data.DeletionKey == deletionKey {
			error := client.Collection.FindOneAndDelete(client.Ctx, bson.M{"slug": slug}).Decode(&url)

			return url, error
		} else {

			return url, errors.New("invalid deletion key")
		}
	}

}
