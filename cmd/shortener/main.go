package main

import (
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
)

func main() {

	httprouters.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", httprouters.NewRouter()))

}
