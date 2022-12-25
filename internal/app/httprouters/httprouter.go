package httprouters

import (
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// маршрутизатор запросов
func NewRouter(hn *handlers.Handler) chi.Router {
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
	rout.Delete("/api/user/urls", hn.Delete)
	// маршрут POST "/api/shorten/batch" пакетная выдача коротких ссылок
	rout.Post("/api/shorten/batch", hn.PutBatch)
	// маршрут GET "/ping" проверка доступности PostgreSQL
	rout.Get("/ping", hn.Ping)
	// маршрут GET "/api/user/urls"  получение ссылок пользователя
	rout.Get("/api/user/urls", hn.GetBatch)
	// маршрут GET "/{id}" получение ссылки по котороткой ссылке
	rout.Get("/{id}", hn.Get)
	// маршрут POST "/api/shorten" выдача короткой ссылки по JSON в теле запроса
	rout.Post("/api/shorten", hn.PutJSON)
	// маршрут POST "/" выдача короткой ссылки по текстовыму URL в теле запроса
	rout.Post("/", hn.Put)
	// возврат ошибки 400 для всех остальных запросов
	rout.HandleFunc("/*", hn.IncorrectRequest)

	return rout
}
