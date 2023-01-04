// Package httprouters - пакет роутера HTTP запросов.
package httprouters

import (
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter маршрутизатор запросов.
func NewRouter(
	putHandler *handlers.PutHandler,
	getHandler *handlers.GetHandler,
	deleteHandler *handlers.DeleteHandler,
	pingHandler *handlers.PingHandler) chi.Router {

	// chi роутер
	rout := chi.NewRouter()
	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	rout.Use(middleware.Logger)
	rout.Use(middleware.Recoverer)
	// дополнительный middleware
	rout.Use(middlewareGzip)
	rout.Use(middlewareCookie)
	// профилировщик
	rout.Mount("/debug", middleware.Profiler())

	// маршрут DELETE "/api/user/urls" пакетное удаление коротки ссылок
	rout.Delete("/api/user/urls", deleteHandler.Delete)
	// маршрут POST "/api/shorten/batch" пакетная выдача коротких ссылок
	rout.Post("/api/shorten/batch", putHandler.PutBatch)
	// маршрут GET "/ping" проверка доступности PostgreSQL
	rout.Get("/ping", pingHandler.Ping)
	// маршрут GET "/api/user/urls"  получение ссылок пользователя
	rout.Get("/api/user/urls", getHandler.GetBatch)
	// маршрут GET "/{id}" получение ссылки по котороткой ссылке
	rout.Get("/{id}", getHandler.Get)
	// маршрут POST "/api/shorten" выдача короткой ссылки по JSON в теле запроса
	rout.Post("/api/shorten", putHandler.PutJSON)
	// маршрут POST "/" выдача короткой ссылки по текстовыму URL в теле запроса
	rout.Post("/", putHandler.Put)

	// возврат ошибки 404 для всех остальных запросов - роутер chi

	return rout
}
