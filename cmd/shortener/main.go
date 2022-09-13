package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/services"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

func main() {
	port := ":8080"
	fmt.Printf("Started server on port %s\n", port)

	s := storage.NewMapStorage("map")
	srvs := services.NewService(s)
	h := handlers.NewHandler(srvs)
	r:= httprouters.NewRouter(h)

	log.Fatal(http.ListenAndServe(port, r))
}
