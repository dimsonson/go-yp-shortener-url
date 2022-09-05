package httprouters

import (
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HTTPRouter(w http.ResponseWriter, r *http.Request) {
	// определяем роутер chi
	rout := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	//rout.Use(middleware.RequestID)
	//rout.Use(middleware.RealIP)
	rout.Use(middleware.Logger)
	rout.Use(middleware.Recoverer)

	rout.Get("/*", handlers.GetHandler)
	rout.Post("/*", handlers.PostHandler)
	rout.HandleFunc("/*", handlers.DefHandler)
}
