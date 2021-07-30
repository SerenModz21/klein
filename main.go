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

	db, error := database.CreateClient(ctx, configuration.Mongo)
	if error != nil {
		log.Panic(error)
	}

	urls := services.UrlClient{
		Ctx:        ctx,
		Collection: db.Collection("urls"),
	}

	router := mux.NewRouter()

	router.HandleFunc("/shorten", routes.ShortenUrl(&urls)).Methods(http.MethodPost)
	router.HandleFunc("/{slug}", routes.RedirectUrl(&urls)).Methods(http.MethodGet)
	router.HandleFunc("/{slug}", routes.DeleteUrl(&urls)).Methods(http.MethodDelete)

	log.Info("now listening on port", configuration.Port)
	http.ListenAndServe(configuration.Port, router)
}
