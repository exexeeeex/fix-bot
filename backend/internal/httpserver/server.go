package httpserver

import (
	"net/http"
)

const allowedOrigin = "http://localhost:3000"

type Handler interface {
	RegisterRoutes(mux *http.ServeMux)
}

func New(handlers ...Handler) *http.Server {
	mux := http.NewServeMux()

	for _, handler := range handlers {
		handler.RegisterRoutes(mux)
	}

	return &http.Server{
		Handler: cors(mux),
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != allowedOrigin {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Vary", "Origin")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
