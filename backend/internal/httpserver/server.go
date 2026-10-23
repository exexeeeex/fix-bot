package httpserver

import (
	"net/http"
)

type Handler interface {
	RegisterRoutes(mux *http.ServeMux)
}

func New(handlers ...Handler) *http.Server {
	mux := http.NewServeMux()

	for _, handler := range handlers {
		handler.RegisterRoutes(mux)
	}

	return &http.Server{
		Handler: mux,
	}
}
