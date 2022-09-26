package httprouters

import (
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var defaultCompressibleContentTypes = []string{
	"text/html",
	"text/css",
	"text/plain",
	"text/javascript",
	"application/javascript",
	"application/x-javascript",
	"application/x-gzip",
	"application/json",
	"application/atom+xml",
	"application/rss+xml",
	"image/svg+xml",
}

func NewRouter(hn *handlers.Handler) chi.Router {
	// chi роутер
	rout := chi.NewRouter()
	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	rout.Use(middleware.Logger)
	rout.Use(middleware.Recoverer)
	//rout.Use(middleware.Compress(5))
	//rout.Use(middleware.Def)
	compressor := middleware.NewCompressor(5, "application/x-gzip")

	rout.Use(compressor.Handler)

	// маршрут GET "/{id}" id в URL
	rout.Get("/{id}", hn.HandlerGetShortURL)
	// маршрут POST "/api/shorten" c JSON в теле запроса
	rout.Post("/api/shorten", hn.HandlerCreateShortJSON)
	// маршрут POST "/" с текстовым URL в теле запроса
	rout.Post("/", hn.HandlerCreateShortURL)
	// возврат ошибки 400 для всех остальных запросов
	rout.HandleFunc("/*", hn.IncorrectRequests)

	return rout
}
