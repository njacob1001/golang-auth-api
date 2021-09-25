package handler

import (
	"net/http"
)

func HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("OK"))

		return
	}
}
