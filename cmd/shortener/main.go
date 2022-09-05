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

	httprouters.HTTPRouter()
	log.Fatal(http.ListenAndServe(":8080", httprouters.NewRouter()))
	
}

