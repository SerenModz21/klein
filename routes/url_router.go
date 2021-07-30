package routes

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"sach.demiboy.me/common"
	"sach.demiboy.me/database/models"
	"sach.demiboy.me/database/services"
)

type NormalResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ShortenResponse struct {
	Success bool   `json:"success"`
	Slug    string `json:"slug"`
	Long    string `json:"long"`
	Short   string `json:"short"`
	Key     string `json:"key"`
}

type DeleteResponse struct {
	Success bool       `json:"success"`
	Url     models.Url `json:"url"`
}

func RedirectUrl(service services.IUrl) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if result, error := service.Get(mux.Vars(r)["slug"]); error != nil {
			rw.WriteHeader(400)
			rw.Write([]byte("<h1>Invalid slug provided.</h1>"))
		} else {
			http.Redirect(rw, r, result.Long, http.StatusMovedPermanently)
		}
	}
}

func ShortenUrl(service services.IUrl) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("url")

		if _, error := url.ParseRequestURI(query); error != nil {
			common.WriteJson(rw, http.StatusBadRequest, NormalResponse{
				Message: "No url query provided or invalid url query provided, example: 'https://link.shortener/shorten?url=https://google.com'",
				Success: false,
			})

			return
		}

		if response, error := service.Insert(models.Url{
			Long:        query,
			Slug:        common.RandomString(4),
			DeletionKey: common.RandomString(8),
		}); error != nil {
			common.WriteJson(rw, http.StatusInternalServerError, NormalResponse{
				Message: error.Error(),
				Success: false,
			})
		} else {
			log.Info(response.Long, " -> ", response.Slug)
			common.WriteJson(rw, http.StatusAccepted, ShortenResponse{
				Success: true,
				Slug:    response.Slug,
				Long:    response.Long,
				Short:   fmt.Sprintf("http://localhost:8080/%s", response.Slug),
				Key:     response.DeletionKey,
			})
		}
	}
}

func DeleteUrl(service services.IUrl) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if data, error := service.Delete(mux.Vars(r)["slug"], r.URL.Query().Get("delete")); error != nil {
			common.WriteJson(rw, http.StatusForbidden, NormalResponse{
				Success: false,
				Message: error.Error(),
			})
		} else {
			common.WriteJson(rw, http.StatusAccepted, DeleteResponse{
				Success: true,
				Url:     data,
			})
		}
	}
}
