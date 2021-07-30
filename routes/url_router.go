package routes

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"sach.demiboy.me/database/models"
	"sach.demiboy.me/database/services"
	"sach.demiboy.me/util"
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
}

func RedirectUrl(service services.IUrl) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		if result, error := service.Get(params["slug"]); error != nil {
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
			util.WriteJson(rw, http.StatusBadRequest, NormalResponse{
				Message: "No url query provided or invalid url query provided, example: 'https://link.shortener/shorten?url=https://google.com'",
				Success: false,
			})

			return
		}

		if response, error := service.Insert(models.Url{
			Long: query,
			Slug: util.RandomString(),
		}); error != nil {
			util.WriteJson(rw, http.StatusInternalServerError, NormalResponse{
				Message: error.Error(),
				Success: false,
			})
		} else {
			util.WriteJson(rw, http.StatusAccepted, ShortenResponse{
				Success: true,
				Slug:    response.Slug,
				Long:    response.Long,
				Short:   fmt.Sprintf("http://localhost:8080/%s", response.Slug),
			})
		}
	}
}
