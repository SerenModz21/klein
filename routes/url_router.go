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

type ResponseData struct {
	Slug  string `json:"slug"`
	Long  string `json:"long"`
	Short string `json:"short"`
	Key   string `json:"key"`
}

type NormalResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type DataResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    ResponseData `json:"data"`
}

type IUrlRouter interface {
	RedirectUrl() http.HandlerFunc
	ShortenUrl() http.HandlerFunc
	DeleteUrl() http.HandlerFunc
}

type URLRouter struct {
	Service services.IUrl
}

func (router URLRouter) RedirectUrl() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if result, error := router.Service.Get(mux.Vars(r)["slug"]); error != nil {
			rw.WriteHeader(400)
			rw.Write([]byte("<h1>Invalid slug provided.</h1>"))
		} else {
			http.Redirect(rw, r, result.Long, http.StatusMovedPermanently)
		}
	}
}

func (router URLRouter) ShortenUrl() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("url")

		if _, error := url.ParseRequestURI(query); error != nil {
			common.WriteJson(rw, http.StatusBadRequest, NormalResponse{
				Message: "No url query provided or invalid url query provided, example: 'https://link.shortener/shorten?url=https://google.com'",
				Success: false,
			})

			return
		}

		if response, error := router.Service.Insert(models.Url{
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
			common.WriteJson(rw, http.StatusAccepted, DataResponse{
				Success: true,
				Message: fmt.Sprintf("Successfully shortened url %s to %s", response.Long, response.Slug),
				Data: ResponseData{
					Slug:  response.Slug,
					Long:  response.Long,
					Short: fmt.Sprintf("http://localhost:8080/%s", response.Slug),
					Key:   response.DeletionKey},
			})
		}
	}
}

func (router URLRouter) DeleteUrl() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if data, error := router.Service.Delete(mux.Vars(r)["slug"], r.URL.Query().Get("delete")); error != nil {
			common.WriteJson(rw, http.StatusForbidden, NormalResponse{
				Success: false,
				Message: fmt.Sprintf("Either the user was not found or you provided an invalid key, error for debugging purposes: %s", error.Error()),
			})
		} else {
			common.WriteJson(rw, http.StatusAccepted, DataResponse{
				Success: true,
				Message: fmt.Sprintf("Successfully deleted slug %s.", data.Slug),
				Data: ResponseData{
					Slug:  data.Slug,
					Long:  data.Long,
					Short: fmt.Sprintf("http://localhost:8080/%s", data.Slug),
					Key:   data.DeletionKey,
				},
			})
		}
	}
}
