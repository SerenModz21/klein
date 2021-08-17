package routes

import (
	"net/http"

	"github.com/aymerick/raymond"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sach.demiboy.me/common"
)

func (router AppRouter) Index() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		names := []string{}

		temp := int64(10)

		if urls, error := router.Service.Search(nil, options.FindOptions{Limit: &temp, Sort: bson.M{"_id": -1}}); error == nil {
			for _, url := range urls {
				names = append(names, url.Slug)
			}
		}

		common.WriteHTML(rw, http.StatusAccepted, raymond.MustRender(router.Templates["index"], map[string]interface{}{
			"version": router.Version,
			"urls":    names,
		}))
	}
}
