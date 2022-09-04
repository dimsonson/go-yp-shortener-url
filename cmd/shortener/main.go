package main

import (
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
)

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", httprouters.HTTPRouter)
	// конструируем сервер
	log.Fatal(http.ListenAndServe(":8080", nil))
}
