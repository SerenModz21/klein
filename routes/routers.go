package routes

import (
	"net/http"

	"sach.demiboy.me/database/services"
)

type IUrlRouter interface {
	RedirectUrl() http.HandlerFunc
	ShortenUrl() http.HandlerFunc
	DeleteUrl() http.HandlerFunc
}

type UrlRouter struct {
	Service services.IUrl
}

type IAppRouter interface {
	Index() http.HandlerFunc
}

type AppRouter struct {
	Templates map[string]string
	Service   services.IUrl
	Version   string
}
