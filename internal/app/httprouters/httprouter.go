package httprouters

import (
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(hn *handlers.Handler) chi.Router {
	rout := chi.NewRouter()
	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	rout.Use(middleware.Logger)
	rout.Use(middleware.Recoverer)
	// маршруты
	rout.Get("/{id}", hn.HandlerGetShortURL)
	rout.Post("/", hn.HandlerCreateShortURL)
	rout.HandleFunc("/*", hn.IncorrectRequests)
	return rout
}
