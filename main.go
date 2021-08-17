package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"sach.demiboy.me/config"
	"sach.demiboy.me/database"
	"sach.demiboy.me/database/services"
	"sach.demiboy.me/routes"
)

func main() {
	version := "v1"

	configuration := config.GetConfig()
	ctx := context.TODO()

	// It's gonna panic and exit either way, don't bother checking for errors.
	db, _ := database.ConnectMongo(ctx, configuration.Mongo)
	redisClient, _ := database.ConnectRedis(ctx, configuration.Redis)

	urlService := services.UrlClient{
		Ctx:        ctx,
		Collection: db.Collection("urls"),
		Redis:      redisClient,
	}

	router := mux.NewRouter()
	urlRouter := routes.UrlRouter{Service: &urlService}

	index, _ := ioutil.ReadFile("views/index.hbs")

	appRouter := routes.AppRouter{
		Version: version,
		Templates: map[string]string{
			"index": string(index),
		},
		Service: &urlService,
	}

	appSubRouter := router.PathPrefix("/").Subrouter()
	appSubRouter.HandleFunc("/", appRouter.Index()).Methods(http.MethodGet)
	appSubRouter.HandleFunc("/{slug}", urlRouter.RedirectUrl()).Methods(http.MethodGet)

	apiSubRouter := router.PathPrefix(fmt.Sprintf("/api/%s", version)).Subrouter()
	apiSubRouter.HandleFunc("/{slug}", urlRouter.DeleteUrl()).Methods(http.MethodDelete)
	apiSubRouter.HandleFunc("/shorten", urlRouter.ShortenUrl()).Methods(http.MethodPost)

	log.Info("now listening on port", configuration.Port)
	http.ListenAndServe(configuration.Port, router)
}
