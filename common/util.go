package common

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
)

func RandomString(len int) string {
	bytes := make([]byte, len)
	rand.Read(bytes)

	return fmt.Sprintf("%x", bytes)
}

func WriteJson(rw http.ResponseWriter, statusCode int, data interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(statusCode)

	if err := json.NewEncoder(rw).Encode(data); err != nil {
		return
	}
}

func WriteHTML(rw http.ResponseWriter, statusCode int, html string) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(statusCode)
	rw.Write([]byte(html))
}
