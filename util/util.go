package util

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
)

func RandomString() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)

	return fmt.Sprintf("%x", bytes)
}

func WriteJson(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return
	}
}
