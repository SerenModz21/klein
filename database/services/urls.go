package services

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sach.demiboy.me/database/models"
)

type IUrl interface {
	Insert(models.Url) (models.Url, error)
	Get(string) (models.Url, error)
}

type UrlClient struct {
	Ctx        context.Context
	Collection *mongo.Collection
}

func (client *UrlClient) Insert(document models.Url) (models.Url, error) {
	url := models.Url{}
	_, error := client.Collection.InsertOne(client.Ctx, document)
	if error != nil {
		return url, error
	}

	return client.Get(document.Slug)
}

func (client *UrlClient) Get(slug string) (models.Url, error) {
	url := models.Url{}

	error := client.Collection.FindOne(client.Ctx, bson.M{"slug": slug}).Decode(&url)

	return url, error
}
