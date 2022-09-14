package httprouters

import (
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
	rout := chi.NewRouter()
	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	rout.Use(middleware.Logger)
	rout.Use(middleware.Recoverer)
	// маршруты
	rout.Get("/{id}", handlers.GetShortURL)
	rout.Post("/", handlers.CreateShortURL)
	rout.HandleFunc("/*", handlers.IncorrectRequests)

	return rout
}
