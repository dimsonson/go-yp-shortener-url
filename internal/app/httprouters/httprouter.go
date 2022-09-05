package httprouters

import (
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
	rout := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	//rout.Use(middleware.RequestID)
	//rout.Use(middleware.RealIP)
	rout.Use(middleware.Logger)
	rout.Use(middleware.Recoverer)

	rout.HandleFunc("/*", handlers.DefHandler)
	rout.Get("/*", handlers.GetHandler)
	rout.Post("/*", handlers.PostHandler)
	//log.Fatal(http.ListenAndServe(":8080", rout))
	return rout
}
