package main

import (
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
)

func main() {
	addr := ":8080"

	log.Printf("Starting server on port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, httprouters.NewRouter()))
}
