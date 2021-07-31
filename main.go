package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"sach.demiboy.me/config"
	"sach.demiboy.me/database"
	"sach.demiboy.me/database/services"
	"sach.demiboy.me/routes"
)

func main() {
	configuration := config.GetConfig()
	ctx := context.TODO()

	// It's gonna panic and exit either way, don't bother checking for errors.
	db, _ := database.ConnectMongo(ctx, configuration.Mongo)
	redisClient, _ := database.ConnectRedis(ctx, configuration.Redis)

	router := mux.NewRouter()
	urlRouter := routes.URLRouter{Service: &services.UrlClient{
		Ctx:        ctx,
		Collection: db.Collection("urls"),
		Redis:      redisClient,
	}}

	router.HandleFunc("/shorten", urlRouter.ShortenUrl()).Methods(http.MethodPost)
	router.HandleFunc("/{slug}", urlRouter.RedirectUrl()).Methods(http.MethodGet)
	router.HandleFunc("/{slug}", urlRouter.DeleteUrl()).Methods(http.MethodDelete)

	log.Info("now listening on port", configuration.Port)
	http.ListenAndServe(configuration.Port, router)
}
