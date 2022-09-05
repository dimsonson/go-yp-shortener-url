package main

import (
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	NewRouter()
	log.Fatal(http.ListenAndServe(":8080", httprouters.NewRouter()))
	
}

func NewRouter() chi.Router { //chi.Router {
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

//log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(NewRouter)))
//log.Fatal(http.ListenAndServe(":8080", rout))